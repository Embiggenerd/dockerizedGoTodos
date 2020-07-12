package models

import (
	"goTodos/utils"
	"net/http"
)

// User is the shape of user Data
type User struct {
	ID        int
	Age       int
	FirstName string
	LastName  string
	Email     string
	Password  string
}

// RegisterUser returns
func RegisterUser(u *User) (*User, *utils.HTTPError) {
	var id int
	sqlUser := `
		INSERT INTO users ( age, first_name, last_name, email, password)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id; `
	err := db.QueryRow(
		sqlUser, u.Age, u.FirstName, u.LastName, u.Email, u.Password).Scan(&id)

	if err != nil {
		return nil, &utils.HTTPError{Code: 400, Msg: "User already registered"}
	}

	user := new(User)

	sqlQuery := `
		SELECT * FROM users WHERE id = $1;`

	row := db.QueryRow(sqlQuery, id)

	err = row.Scan(&user.ID, &user.Age, &user.FirstName,
		&user.LastName, &user.Email, &user.Password)

	if err != nil {
		return nil, &utils.HTTPError{
			Code: http.StatusInternalServerError,
			Msg:  "An error happend that wasn't your fault"}
	}

	return user, nil
}

func validatePassword(storedPassword, providedPassword string) (bool, *utils.HTTPError) {
	if storedPassword != providedPassword {
		return false, &utils.HTTPError{
			Code: http.StatusUnauthorized,
			Msg:  "Try again",
		}
	}
	return true, nil
}

// LoginUser compares submitted password to the one stored in users
// table, returns user if validated
func LoginUser(p, e string) (*User, *utils.HTTPError) {
	sqlEmailQuery := `
		SELECT * FROM users WHERE EMAIL = $1
		LIMIT 1;`
	user := new(User)

	err := db.QueryRow(sqlEmailQuery, e).Scan(&user.ID, &user.Age,
		&user.FirstName, &user.LastName, &user.Email, &user.Password)

	if err != nil {
		return nil, &utils.HTTPError{
			Code: http.StatusBadRequest,
			Msg:  "Email not found in database",
		}
	}

	isValid, httpErr := validatePassword(user.Password, p)
	if err != nil && isValid != true {
		return nil, httpErr
	}
	return user, nil
}

// GetUserFromSession retrieves a user from a cookie by looking up a randomly generated
// hex value on a cookie.
func GetUserFromSession(hex string) (*User, *utils.HTTPError) {
	var id int
	sqlSessionQuery := `
		 SELECT userid FROM sessions
		 WHERE hex = $1`
	user := new(User)
	err := db.QueryRow(sqlSessionQuery, hex).Scan(&id)

	sqlUserQuery := `
		SELECT * FROM users
		WHERE id = $1;`
	err = db.QueryRow(sqlUserQuery, id).Scan(&user.ID, &user.Age,
		&user.FirstName, &user.LastName, &user.Email, &user.Password)

	if err != nil {
		return nil, &utils.HTTPError{
			Code: http.StatusInternalServerError,
			Msg:  "There was an error that isn't your fault",
		}
	}

	return user, nil
}

// DeleteSession deletes old session on login for space efficiency
func DeleteSession(userID int) *utils.HTTPError {
	sqlSession := `
		DELETE FROM sessions
		WHERE userid = $1;`

	_, err := db.Query(sqlSession, userID)
	if err != nil {
		return &utils.HTTPError{
			Code: http.StatusInternalServerError,
			Msg:  "There was an error that isn't your fault",
		}
	}
	return nil
}

// CreateSession inserts user id, random hex value for
// Fetching user for auth
func CreateSession(hex string, userID int) *utils.HTTPError {
	sqlSession := `
		INSERT INTO sessions ( hex, userid )
		VALUES( $1, $2);`

	_, err := db.Query(sqlSession, hex, userID)
	if err != nil {
		return &utils.HTTPError{
			Code: http.StatusInternalServerError,
			Msg:  "There was an error that isn't your fault",
		}
	}
	return nil
}

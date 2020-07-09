package models

// DeleteSession deletes old session on login for space efficiency
func DeleteSession(userID int) error {
	sqlSession := `
		DELETE FROM sessions
		WHERE userid = $1;`

	_, err := db.Query(sqlSession, userID)
	if err != nil {
		return err
	}
	return nil
}

// CreateSession inserts user id, random hex value for
// Fetching user for auth
func CreateSession(hex string, userID int) error {
	sqlSession := `
		INSERT INTO sessions ( hex, userid )
		VALUES( $1, $2);`

	_, err := db.Query(sqlSession, hex, userID)
	if err != nil {
		return err
	}
	return nil
}

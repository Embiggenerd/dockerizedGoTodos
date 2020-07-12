package models

import (
	"goTodos/utils"
	"log"
	"net/http"
)

// Todo is for storing values returned from query
type Todo struct {
	ID       int
	Body     string
	AuthorID int
	Done     bool
}

// EditTodo takes id, new body, updates todo table with id to have new body
func EditTodo(id, body string) (*Todo, *utils.HTTPError) {
	sqlUpdate := `
		UPDATE todos
		SET body = $1
		WHERE id = $2`
	_, err := db.Exec(sqlUpdate, body, id)

	todo := new(Todo)
	sqlQuery := `SELECT * FROM todos WHERE id = $1`

	row := db.QueryRow(sqlQuery, id)
	err = row.Scan(&todo.ID, &todo.Body, &todo.AuthorID, &todo.Done)

	if err != nil {
		return nil, &utils.HTTPError{
			Code: http.StatusInternalServerError,
			Msg:  "There was an error that isn't your fault",
		}
	}

	return todo, nil
}

// DeleteTodo sends delete instuction to db with todo's id
func DeleteTodo(id string) *utils.HTTPError {
	sqlDelete := `
		DELETE FROM todos
		WHERE id = $1`

	_, err := db.Exec(sqlDelete, id)

	if err != nil {
		return &utils.HTTPError{
			Code: http.StatusInternalServerError,
			Msg:  "There was an error that isn't your fault",
		}
	}
	return nil
}

// SubmitTodo inserts values into todo table, querys by returned id, returns added todo
func SubmitTodo(t *Todo) (*Todo, *utils.HTTPError) {
	id := 0

	sqlInsert := `
		INSERT INTO todos ( body, authorId, done)
		VALUES ($1, $2, $3)
		RETURNING id`

	err := db.QueryRow(sqlInsert, t.Body, t.AuthorID, t.Done).Scan(&id)

	todo := new(Todo)

	sqlQuery := `SELECT * FROM todos WHERE id = $1`
	row := db.QueryRow(sqlQuery, id)

	err = row.Scan(&todo.ID, &todo.Body, &todo.AuthorID, &todo.Done)

	if err != nil {
		return nil, &utils.HTTPError{
			Code: http.StatusInternalServerError,
			Msg:  "There was an error that isn't your fault",
		}
	}
	return todo, nil
}

// GetTodos returns all todos in database
func GetTodos(userID int) ([]*Todo, *utils.HTTPError) {
	rows, err := db.Query(`
		SELECT * FROM todos
		WHERE authorid = $1;`, userID)

	defer rows.Close()

	todos := make([]*Todo, 0)

	for rows.Next() {
		todo := new(Todo)
		err := rows.Scan(&todo.ID, &todo.Body, &todo.AuthorID, &todo.Done)
		if err != nil {
			log.Fatal(err)
		}
		todos = append(todos, todo)
		log.Println(todo)
	}

	err = rows.Err()

	if err != nil {
		return nil, &utils.HTTPError{
			Code: http.StatusInternalServerError,
			Msg:  "There was an error that isn't your fault",
		}
	}

	return todos, nil
}

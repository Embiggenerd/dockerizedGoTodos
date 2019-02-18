package routes

import (
	"goTodos/models"
	"goTodos/utils"
	"net/http"
)

func editHandler(w http.ResponseWriter, r *http.Request) {
	todoId := r.URL.Path[len("/edit/"):]

	if r.Method == "GET" {
		err := tmplts.ExecuteTemplate(w, "index.html",
			templData{
				State:  "editTodo",
				Header: "Edit your todo",
				Styles: cacheBustedCss,
				TodoId: todoId,
				Todos:  nil,
				User:   nil,
			})

		if err != nil {
			utils.InternalServerError(w, r)
		}

	} else {
		r.ParseForm()
		body := r.Form["body"][0]

		_, err := models.EditTodo(todoId, body)

		if err != nil {
			utils.InternalServerError(w, r)
		}
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

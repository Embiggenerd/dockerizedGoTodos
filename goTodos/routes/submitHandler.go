package routes

import (
	"goTodos/models"
	"goTodos/utils"
	"net/http"
)

func submitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		err := tmplts.ExecuteTemplate(w, "index.html", templData{
			State:  "submitTodo",
			Header: "Submit a new todo",
			Styles: cacheBustedCss,
			TodoId: "",
			Todos:  nil,
			User:   nil,
		})

		if err != nil {
			utils.InternalServerError(w, r)
		}

	} else {
		user, ok := r.Context().Value(contextKey("user")).(*models.User)

		if !ok {
			utils.InternalServerError(w, r)
		}

		r.ParseForm()
		todo := models.Todo{
			ID:       0,
			Body:     r.Form["body"][0],
			AuthorID: user.ID,
			Done:     false,
		}

		_, err := models.SubmitTodo(&todo)

		if err != nil {
			utils.InternalServerError(w, r)
		}

		http.Redirect(w, r, "/", http.StatusFound)
	}
}

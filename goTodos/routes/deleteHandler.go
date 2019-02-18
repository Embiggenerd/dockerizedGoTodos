package routes

import (
	"goTodos/models"
	"goTodos/utils"
	"net/http"
)

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id := r.URL.Path[len("/delete/"):]
		err := models.DeleteTodo(id)
		if err != nil {
			utils.InternalServerError(w, r)
		}
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

package routes

import (
	"goTodos/models"
	"goTodos/utils"
	"net/http"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// Get user off context
	user, ok := r.Context().Value(contextKey("user")).(*models.User)
	// Any error is 500 status
	if !ok {
		utils.InternalServerError(w, r)
	}
	// Get todos from todos model
	todos, err := models.GetTodos(user.ID)
	// Also 500 error
	if err != nil {
		utils.InternalServerError(w, r)
	}
	// Render template
	err = tmplts.ExecuteTemplate(w, "index.html",
		templData{
			State:  "home",
			Header: "Home",
			Styles: cacheBustedCss,
			TodoId: "",
			Todos:  todos,
			User:   user,
		})
	if err != nil {
		utils.InternalServerError(w, r)
	}
}

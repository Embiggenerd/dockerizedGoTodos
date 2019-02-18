package routes

import (
	"goTodos/models"
	"goTodos/utils"
	"net/http"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(contextKey("user")).(*models.User)
	if !ok {
		utils.InternalServerError(w, r)
	}

	todos, err := models.GetTodos(user.ID)
	if err != nil {
		utils.InternalServerError(w, r)
	}

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

package routes

import (
	"fmt"
	"goTodos/models"
	"goTodos/utils"
	"net/http"
	"strconv"
)

func registerUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		err := tmplts.ExecuteTemplate(w, "index.html", templData{
			State: "signup", Header: "Register with an email and password", Styles: cacheBustedCss, TodoId: "", Todos: nil, User: nil,
		})

		if err != nil {
			utils.InternalServerError(w, r)
		}

	} else {
		r.ParseForm()
		age, err := strconv.Atoi(r.Form["age"][0])

		if err != nil {
			fmt.Println(err)
			utils.InternalServerError(w, r)
		}

		user := models.User{
			ID:        0,
			Age:       age,
			FirstName: r.Form["firstName"][0],
			LastName:  r.Form["lastName"][0],
			Email:     r.Form["email"][0],
			Password:  r.Form["password"][0]}

		_, err = models.RegisterUser(&user)

		if err != nil {
			fmt.Println(err)
			utils.InternalServerError(w, r)
		}

		http.Redirect(w, r, "/", http.StatusFound)
	}
}

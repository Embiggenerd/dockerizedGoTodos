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

		tmplts.ExecuteTemplate(w, "index.html", templData{
			State: "signup", Header: "Register with an email and password, all values are required", Styles: cacheBustedCss, TodoId: "", Todos: nil, User: nil, Error: nil,
		})

	} else {

		r.ParseForm()
		age, err := strconv.Atoi(r.Form["age"][0])

		if err != nil {
			fmt.Println(err)
		}

		user := models.User{
			ID:        0,
			Age:       age,
			FirstName: r.Form["firstName"][0],
			LastName:  r.Form["lastName"][0],
			Email:     r.Form["email"][0],
			Password:  r.Form["password"][0]}

		_, httpErr := models.RegisterUser(&user)

		if httpErr != nil {
			w.WriteHeader(httpErr.Code)
			tmplts.ExecuteTemplate(w, "index.html", templData{
				State:  "signup",
				Header: "Register with an email and password, all values are required",
				Styles: cacheBustedCss,
				TodoId: "",
				Todos:  nil,
				User:   nil,
				Error:  &utils.HTTPError{Code: httpErr.Code, Name: httpErr.Name, Msg: httpErr.Msg},
			})
			return
		}

		http.Redirect(w, r, "/login", http.StatusFound)
	}
}

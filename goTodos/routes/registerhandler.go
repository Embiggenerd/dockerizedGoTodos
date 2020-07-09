package routes

import (
	"fmt"
	"goTodos/models"
	"goTodos/utils"
	"net/http"
	"regexp"
	"strconv"
)

func registerUserHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {

		tmplts.ExecuteTemplate(w, "index.html", templData{
			State: "signup", Header: "Register with an email and password", Styles: cacheBustedCss, TodoId: "", Todos: nil, User: nil, Error: nil,
		})

	} else {
		var age int
		var err error
		r.ParseForm()

		fmt.Println("len(age)", len(r.Form["age"]))

		if r.Form["email"][0] == "" {
			w.WriteHeader(400)
			tmplts.ExecuteTemplate(w, "index.html", templData{
				State:  "signup",
				Header: "Register with an email and password",
				Styles: cacheBustedCss,
				TodoId: "",
				Todos:  nil,
				User:   nil,
				Error: &utils.HTTPError{
					Code: 400,
					Name: "Invalid Input",
					Msg:  "Please include an email",
				},
			})
			return
		} else {
			regex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
			if !regex.MatchString(r.Form["email"][0]) {
				w.WriteHeader(400)
				tmplts.ExecuteTemplate(w, "index.html", templData{
					State:  "signup",
					Header: "Register with an email and password",
					Styles: cacheBustedCss,
					TodoId: "",
					Todos:  nil,
					User:   nil,
					Error: &utils.HTTPError{
						Code: 400,
						Name: "Invalid Input",
						Msg:  "Email must be in the form 'name@example.com",
					},
				})
				return
			}
		}

		if r.Form["password"][0] == "" {
			w.WriteHeader(400)
			tmplts.ExecuteTemplate(w, "index.html", templData{
				State:  "signup",
				Header: "Register with an email and password",
				Styles: cacheBustedCss,
				TodoId: "",
				Todos:  nil,
				User:   nil,
				Error: &utils.HTTPError{
					Code: 400,
					Name: "Invalid Input",
					Msg:  "Please include an password",
				},
			})
			return
		}

		if r.Form["age"][0] != "" {
			age, err = strconv.Atoi(r.Form["age"][0])

			if err != nil {
				fmt.Println(err)
				w.WriteHeader(400)
				tmplts.ExecuteTemplate(w, "index.html", templData{
					State:  "signup",
					Header: "Register with an email and password",
					Styles: cacheBustedCss,
					TodoId: "",
					Todos:  nil,
					User:   nil,
					Error: &utils.HTTPError{
						Code: 400,
						Name: "Invalid Input",
						Msg:  "Age must be an integer",
					},
				})
				return
			}
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

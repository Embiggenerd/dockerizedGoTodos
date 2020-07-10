package routes

import (
	"context"
	"fmt"
	"goTodos/models"
	"goTodos/utils"
	"net/http"
	"regexp"
	"strconv"
)

func authRequired(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var sessionHexFromCookie string

		cookie, err := r.Cookie("user-session")

		if err != nil {
			fmt.Println(err)
			err = tmplts.ExecuteTemplate(w, "index.html", templData{
				State:  "withoutAuth",
				Header: "Welcome to Go Postgres Todos",
				TodoId: "",
				Todos:  nil,
				User:   nil,
			})

			if err != nil {
				utils.InternalServerError(w, r)
			}

		} else {
			sessionHexFromCookie = cookie.Value

			// Geg hex token from cookie to find a user in a session row
			user, err := models.GetUserFromSession(sessionHexFromCookie)
			if err != nil {
				utils.UnauthorizedUserError(w)
			}

			f := func(ctx context.Context, k contextKey) {
				v := ctx.Value(k)
				if v != nil {
					fmt.Println("user value in context", v)
					return
				}
				utils.UnauthorizedUserError(w)
			}
			k := contextKey("user")
			ctx := context.WithValue(context.Background(), k, user)
			f(ctx, k)
			handler.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}

func validRegisterBody(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "POST" {

			r.ParseForm()

			fmt.Println("len(age)", len(r.Form["age"]))

			if r.Form["email"][0] == "" {
				w.WriteHeader(400)
				tmplts.ExecuteTemplate(w, "index.html", templData{
					State:  "signup",
					Header: "Register with an email and password",
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
			}

			regex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
			if !regex.MatchString(r.Form["email"][0]) {
				w.WriteHeader(400)
				tmplts.ExecuteTemplate(w, "index.html", templData{
					State:  "signup",
					Header: "Register with an email and password",
					TodoId: "",
					Todos:  nil,
					User:   nil,
					Error: &utils.HTTPError{
						Code: 400,
						Name: "Invalid Input",
						Msg:  "Email must be in the form 'name@example.com'",
					},
				})
				return
			}

			if r.Form["password"][0] == "" {
				w.WriteHeader(400)
				tmplts.ExecuteTemplate(w, "index.html", templData{
					State:  "signup",
					Header: "Register with an email and password",
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
				_, err := strconv.Atoi(r.Form["age"][0])

				if err != nil {
					fmt.Println(err)
					w.WriteHeader(400)
					tmplts.ExecuteTemplate(w, "index.html", templData{
						State:  "signup",
						Header: "Register with an email and password",
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
		}
		handler.ServeHTTP(w, r)
	}
}

func validLoginBody(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "POST" {
			r.ParseForm()

			if r.Form["email"][0] == "" {
				w.WriteHeader(400)
				tmplts.ExecuteTemplate(w, "index.html", templData{
					State:  "login",
					Header: "Log in with an email and password",
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
			}

			regex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
			if !regex.MatchString(r.Form["email"][0]) {
				w.WriteHeader(400)
				tmplts.ExecuteTemplate(w, "index.html", templData{
					State:  "login",
					Header: "Log in with an email and password",
					TodoId: "",
					Todos:  nil,
					User:   nil,
					Error: &utils.HTTPError{
						Code: 400,
						Name: "Invalid Input",
						Msg:  "Email must be in the form 'name@example.com'",
					},
				})
				return
			}

			if r.Form["password"][0] == "" {
				w.WriteHeader(400)
				tmplts.ExecuteTemplate(w, "index.html", templData{
					State:  "login",
					Header: "Log in with an email and password",
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
		}
		handler.ServeHTTP(w, r)
	}
}

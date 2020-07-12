package routes

import (
	"context"
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
			w.WriteHeader(http.StatusOK)

			tmplts.ExecuteTemplate(w, "index.html", templData{
				State:  "withoutAuth",
				Header: "Welcome to Go Postgres Todos",
				TodoId: "",
				Todos:  nil,
				User:   nil,
			})

		} else {
			sessionHexFromCookie = cookie.Value

			// Geg hex token from cookie to find a user in a session row
			user, err := models.GetUserFromSession(sessionHexFromCookie)

			if err != nil {
				respondWithError(w, "withoutAuth", err)
				return
			}

			f := func(ctx context.Context, k contextKey) {
				v := ctx.Value(k)
				if v != nil {
					return
				}
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

			if r.Form["email"][0] == "" {

				respondWithError(
					w,
					"signup",
					&utils.HTTPError{
						Code: http.StatusUnprocessableEntity,
						Msg:  "Invalid Input: Please include an email",
					},
				)

				return
			}

			regex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

			if !regex.MatchString(r.Form["email"][0]) {

				respondWithError(
					w,
					"signup",
					&utils.HTTPError{
						Code: http.StatusUnprocessableEntity,
						Msg:  "Email must be in the form 'name@example.com'",
					},
				)

				return
			}

			if r.Form["password"][0] == "" {

				respondWithError(
					w,
					"signup",
					&utils.HTTPError{
						Code: http.StatusUnprocessableEntity,
						Msg:  "Please include a password",
					},
				)

				return
			}

			if r.Form["age"][0] != "" {
				_, err := strconv.Atoi(r.Form["age"][0])

				if err != nil {

					respondWithError(
						w,
						"signup",
						&utils.HTTPError{
							Code: http.StatusUnprocessableEntity,
							Msg:  "Age must be an integer",
						},
					)
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

				respondWithError(
					w,
					"login",
					&utils.HTTPError{
						Code: http.StatusUnprocessableEntity,
						Msg:  "Please include an email",
					},
				)

				return
			}

			regex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

			if !regex.MatchString(r.Form["email"][0]) {

				respondWithError(
					w,
					"login",
					&utils.HTTPError{
						Code: http.StatusUnprocessableEntity,
						Msg:  "Email must be in the form 'name@example.com'",
					},
				)

				return
			}

			if r.Form["password"][0] == "" {

				respondWithError(
					w,
					"login",
					&utils.HTTPError{
						Code: http.StatusUnprocessableEntity,
						Msg:  "Please include a password",
					},
				)

				return
			}
		}
		handler.ServeHTTP(w, r)
	}
}

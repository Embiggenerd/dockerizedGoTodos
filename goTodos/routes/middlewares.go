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

				errorResp.RespondWithError(
					w,
					http.StatusUnprocessableEntity,
					"Invalid Input: Please include an email",
					"signup",
				)

				return
			}

			regex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

			if !regex.MatchString(r.Form["email"][0]) {

				errorResp.RespondWithError(
					w,
					http.StatusUnprocessableEntity,
					"Email must be in the form 'name@example.com'",
					"signup",
				)

				return
			}

			if r.Form["password"][0] == "" {

				errorResp.RespondWithError(
					w,
					http.StatusUnprocessableEntity,
					"Please include a password",
					"signup",
				)

				return
			}

			if r.Form["age"][0] != "" {
				_, err := strconv.Atoi(r.Form["age"][0])

				if err != nil {

					errorResp.RespondWithError(
						w,
						http.StatusUnprocessableEntity,
						"Age must be an integer",
						"signup",
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

				errorResp.RespondWithError(
					w,
					http.StatusUnprocessableEntity,
					"Please include an email",
					"login",
				)

				return
			}

			regex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

			if !regex.MatchString(r.Form["email"][0]) {
				errorResp.RespondWithError(
					w,
					http.StatusUnprocessableEntity,
					"Email must be in the form 'name@example.com'",
					"login",
				)
				return
			}

			if r.Form["password"][0] == "" {

				errorResp.RespondWithError(
					w,
					http.StatusUnprocessableEntity,
					"Please include an password",
					"login",
				)

				return
			}
		}
		handler.ServeHTTP(w, r)
	}
}

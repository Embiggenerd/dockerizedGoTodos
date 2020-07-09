package routes

import (
	"context"
	"fmt"
	"goTodos/models"
	"goTodos/utils"
	"net/http"
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
				Styles: cacheBustedCss,
				TodoId: "",
				Todos:  nil,
				User:   nil,
			})

			if err != nil {
				utils.InternalServerError(w, r)
			}

		} else {
			sessionHexFromCookie = cookie.Value

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

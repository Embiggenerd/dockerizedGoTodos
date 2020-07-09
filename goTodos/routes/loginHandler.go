package routes

import (
	"goTodos/models"
	"goTodos/utils"
	"net/http"
)

func loginUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		err := tmplts.ExecuteTemplate(w, "index.html", templData{State: "login", Header: "Log in with an email and password", Styles: cacheBustedCss, TodoId: "", Todos: nil, User: nil})

		if err != nil {
			utils.InternalServerError(w, r)
		}

	} else if r.Method == "POST" {
		r.ParseForm()

		user, err := models.LoginUser(r.Form["password"][0], r.Form["email"][0])

		if err != nil {
			w.WriteHeader(400)
			tmplts.ExecuteTemplate(w, "index.html", templData{
				State:  "login",
				Header: "Log in with an email and password",
				Styles: cacheBustedCss,
				TodoId: "",
				Todos:  nil,
				User:   nil,
				Error: &utils.HTTPError{
					Code: 400,
					Name: "Invalid Input",
					Msg:  "Try again",
				},
			})
		} else {
			// Delete old session
			err = models.DeleteSession(user.ID)

			if err != nil {
				utils.InternalServerError(w, r)
			}

			hex, err := utils.RandHex(10)

			if err != nil {
				utils.InternalServerError(w, r)
			}
			// Create new session
			err = models.CreateSession(hex, user.ID)

			if err != nil {
				utils.InternalServerError(w, r)
			}

			cookie := &http.Cookie{
				Name:     "user-session",
				Value:    hex,
				MaxAge:   60 * 60 * 24,
				HttpOnly: true,
			}
			http.SetCookie(w, cookie)
			http.Redirect(w, r, "/", http.StatusFound)
		}
	}
}

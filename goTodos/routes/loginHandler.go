package routes

import (
	"fmt"
	"goTodos/models"
	"goTodos/utils"
	"net/http"
)

// Validate password, if true:
//	Return user data
// 	Find old session by user id, delete
//	Create random hex string
//	Create new row in sessions table with new user id, hex
//
func loginUserHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("LoginHandlerInvoked")
	if r.Method == "GET" {
		err := tmplts.ExecuteTemplate(w, "index.html", templData{State: "login", Header: "Log in with an email and password", Styles: cacheBustedCss, TodoId: "", Todos: nil, User: nil})

		if err != nil {
			utils.InternalServerError(w, r)
		}

	} else if r.Method == "POST" {
		r.ParseForm()

		user, err := models.LoginUser(r.Form["password"][0], r.Form["email"][0])

		if err != nil {
			// http.Redirect(w, r, "/register", http.StatusFound)
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
			err = models.DeleteSession(user.ID)

			if err != nil {
				utils.InternalServerError(w, r)
			}

			hex, err := utils.RandHex(10)

			if err != nil {
				utils.InternalServerError(w, r)
			}

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

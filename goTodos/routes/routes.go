package routes

import (
	_ "expvar"
	"goTodos/models"
	"goTodos/utils"
	"net/http"
	_ "net/http/pprof"
	"strconv"
	"text/template"
)

var tmplts = template.Must(template.ParseFiles("views/index.html", "views/withoutAuth.html", "views/home.html", "views/nav.html",
	"views/head.html", "views/header.html", "views/error.html", "views/footer.html", "views/login.html", "views/editTodo.html", "views/signup.html", "views/submitTodo.html"))

type templData struct {
	State  string
	Header string
	TodoId string
	Todos  []*models.Todo
	User   *models.User
	Error  *utils.HTTPError
}

type contextKey string

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// Get user off context
	user, ok := r.Context().Value(contextKey("user")).(*models.User)
	// Any error is 500 status
	if !ok {
		utils.InternalServerError(w, r)
	}
	// Get todos from todos model
	todos, err := models.GetTodos(user.ID)
	// Also 500 error
	if err != nil {
		utils.InternalServerError(w, r)
	}
	// Render template
	err = tmplts.ExecuteTemplate(w, "index.html",
		templData{
			State:  "home",
			Header: "Home",
			TodoId: "",
			Todos:  todos,
			User:   user,
		})
	if err != nil {
		utils.InternalServerError(w, r)
	}
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		err := tmplts.ExecuteTemplate(w, "index.html", templData{
			State:  "submitTodo",
			Header: "Submit a new todo",
			TodoId: "",
			Todos:  nil,
			User:   nil,
		})

		if err != nil {
			utils.InternalServerError(w, r)
		}

	} else {
		user, ok := r.Context().Value(contextKey("user")).(*models.User)

		if !ok {
			utils.InternalServerError(w, r)
		}

		r.ParseForm()
		todo := models.Todo{
			ID:       0,
			Body:     r.Form["body"][0],
			AuthorID: user.ID,
			Done:     false,
		}

		_, err := models.SubmitTodo(&todo)

		if err != nil {
			utils.InternalServerError(w, r)
		}

		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	todoId := r.URL.Path[len("/edit/"):]

	if r.Method == "GET" {
		err := tmplts.ExecuteTemplate(w, "index.html",
			templData{
				State:  "editTodo",
				Header: "Edit your todo",
				TodoId: todoId,
				Todos:  nil,
				User:   nil,
			})

		if err != nil {
			utils.InternalServerError(w, r)
		}

	} else {
		r.ParseForm()
		body := r.Form["body"][0]

		_, err := models.EditTodo(todoId, body)

		if err != nil {
			utils.InternalServerError(w, r)
		}
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id := r.URL.Path[len("/delete/"):]
		err := models.DeleteTodo(id)
		if err != nil {
			utils.InternalServerError(w, r)
		}
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func registerUserHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {

		tmplts.ExecuteTemplate(w, "index.html", templData{
			State:  "signup",
			Header: "Register with an email and password",
			TodoId: "",
			Todos:  nil,
			User:   nil,
			Error:  nil,
		})

	} else if r.Method == "POST" {

		age, _ := strconv.Atoi(r.Form["age"][0])

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

func loginUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		err := tmplts.ExecuteTemplate(w, "index.html", templData{
			State:  "login",
			Header: "Log in with an email and password",
			TodoId: "",
			Todos:  nil,
			User:   nil})

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

func logoutUserHandler(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:     "user-session",
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", http.StatusFound)
}

// Init initializes routes in main
func Init() {
	fs := http.FileServer(http.Dir("static/"))

	http.Handle("/static/", http.StripPrefix("/static", fs))
	http.HandleFunc("/", authRequired(indexHandler))
	http.HandleFunc("/submit", authRequired(submitHandler))
	http.HandleFunc("/edit/", authRequired(editHandler))
	http.HandleFunc("/delete/", deleteHandler)
	http.HandleFunc("/register", validRegisterBody(registerUserHandler))
	http.HandleFunc("/login", validLoginBody(loginUserHandler))
	http.HandleFunc("/logout", logoutUserHandler)
	// http.HandleFunc("/debug/pprof", pprof.Index)
	// http.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	// http.HandleFunc("/debug/pprof/profile", pprof.Profile)
	// http.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	// http.HandleFunc("/debug/pprof/trace", pprof.Trace)
	// http.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	// http.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	// http.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	// http.Handle("/debug/pprof/block", pprof.Handler("block"))
	// http.Handle("/debug/vars", http.DefaultServeMux)
	http.ListenAndServe(":8080", nil)
}

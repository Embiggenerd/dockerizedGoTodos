package routes

import (
	_ "expvar"
	"goTodos/logger"
	"goTodos/models"
	"goTodos/utils"
	"log"
	"net/http"
	_ "net/http/pprof"
	"strconv"
	"text/template"

	"github.com/improbable-eng/go-httpwares/logging/logrus/ctxlogrus"
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

// If user gets to this, he will be served index with 'home' state,
// which shows him his todos
func indexHandler(w http.ResponseWriter, r *http.Request) {
	ctxlogrus.Extract(r.Context()).Info("logging")

	// Get user off context
	user, _ := r.Context().Value(contextKey("user")).(*models.User)

	// Get todos from todos model
	todos, err := models.GetTodos(user.ID)

	// Also 500 error
	if err != nil {
		respondWithError(
			w,
			"home",
			&utils.HTTPError{
				Code: http.StatusInternalServerError,
				Msg:  "An error happend that isn't your fault",
			},
		)

		return
	}

	// Render template
	w.WriteHeader(http.StatusOK)
	tmplts.ExecuteTemplate(w, "index.html",
		templData{
			State:  "home",
			Header: "Home",
			TodoId: "",
			Todos:  todos,
			User:   user,
			Error:  nil,
		})
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {

		w.WriteHeader(http.StatusOK)
		tmplts.ExecuteTemplate(w, "index.html", templData{
			State:  "submitTodo",
			Header: "Submit a new todo",
			TodoId: "",
			Todos:  nil,
			User:   nil,
		})

	} else if r.Method == "POST" {
		user, _ := r.Context().Value(contextKey("user")).(*models.User)

		r.ParseForm()

		todo := models.Todo{
			ID:       0,
			Body:     r.Form["body"][0],
			AuthorID: user.ID,
			Done:     false,
		}

		_, err := models.SubmitTodo(&todo)

		if err != nil {
			respondWithError(w, "submitTodo", err)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	todoId := r.URL.Path[len("/edit/"):]

	if r.Method == "GET" {
		tmplts.ExecuteTemplate(w, "index.html",
			templData{
				State:  "editTodo",
				Header: "Edit your todo",
				TodoId: todoId,
				Todos:  nil,
				User:   nil,
			})

	} else {
		r.ParseForm()
		body := r.Form["body"][0]

		_, err := models.EditTodo(todoId, body)

		if err != nil {
			respondWithError(w, "editTodo", err)
		}
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id := r.URL.Path[len("/delete/"):]

		err := models.DeleteTodo(id)

		if err != nil {
			respondWithError(w, "home", err)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		w.WriteHeader(http.StatusOK)

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
			Password:  r.Form["password"][0],
		}

		_, err := models.RegisterUser(&user)

		if err != nil {
			respondWithError(w, "signup", err)
			return
		}

		http.Redirect(w, r, "/login", http.StatusFound)
	}
}

func loginUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.WriteHeader(http.StatusOK)

		tmplts.ExecuteTemplate(w, "index.html", templData{
			State:  "login",
			Header: "Log in with an email and password",
			TodoId: "",
			Todos:  nil,
			User:   nil})

	} else if r.Method == "POST" {
		r.ParseForm()

		// Login the user
		user, err := models.LoginUser(r.Form["password"][0], r.Form["email"][0])

		// Delete old session
		err = models.DeleteSession(user.ID)

		// Make a new random hex value
		hex, err := utils.RandHex(10)

		// Create new session
		err = models.CreateSession(hex, user.ID)

		if err != nil {
			respondWithError(w, "login", err)

			return
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
	mux := http.DefaultServeMux
	fs := http.FileServer(http.Dir("static/"))

	mux.Handle("/static/", http.StripPrefix("/static", fs))
	mux.HandleFunc("/", authRequired(indexHandler))
	mux.HandleFunc("/submit", authRequired(submitHandler))
	mux.HandleFunc("/edit/", authRequired(editHandler))
	mux.HandleFunc("/delete/", deleteHandler)
	mux.HandleFunc("/register", validRegisterBody(RegisterUserHandler))
	mux.HandleFunc("/login", validLoginBody(loginUserHandler))
	mux.HandleFunc("/logout", logoutUserHandler)
	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		http.ServeFile(w, r, "favicon.ico")
	})

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: logger.LoggingMiddleware(Metrics(mux)),
	}

	log.Fatal(httpServer.ListenAndServe())
}

package routes

import (
	"goTodos/models"
	"goTodos/utils"
	"net/http"
	"text/template"
)

var tmplts = template.Must(template.ParseFiles("views/index.html", "views/withoutAuth.html", "views/home.html", "views/nav.html",
	"views/head.html", "views/header.html", "views/500.html", "views/footer.html", "views/login.html", "views/editTodo.html", "views/signup.html", "views/submitTodo.html"))

type templData struct {
	State  string
	Header string
	Styles string
	TodoId string
	Todos  []*models.Todo
	User   *models.User
}

type contextKey string

var cacheBustedCss string

// Init initializes routes in main
func Init() {
	// var err string
	cacheBustedCss, _ = utils.BustaCache("mainFloats.css", cacheBustedCss)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	fs := http.FileServer(http.Dir("public/"))
	http.Handle("/static/", http.StripPrefix("/static", fs))
	http.HandleFunc("/", authRequired(indexHandler))
	http.HandleFunc("/submit", authRequired(submitHandler))
	http.HandleFunc("/edit/", authRequired(editHandler))
	http.HandleFunc("/delete/", deleteHandler)
	http.HandleFunc("/register", registerUserHandler)
	http.HandleFunc("/login", loginUserHandler)
	http.HandleFunc("/logout", logoutUserHandler)
	http.HandleFunc("/oops", status500Handler)

	http.ListenAndServe(":8000", nil)
}

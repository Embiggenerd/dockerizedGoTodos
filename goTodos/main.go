package main

import (
	"goTodos/models"
	"goTodos/routes"

	_ "github.com/lib/pq"
)

func main() {
	models.Init()
	routes.Init()
	// var err string
	// cacheBustedCss, err = utils.BustaCache("mainFloats.css", cacheBustedCss)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fs := http.FileServer(http.Dir("public/"))
	// http.Handle("/static/", http.StripPrefix("/static", fs))
	// http.HandleFunc("/", authRequired(indexHandler))
	// http.HandleFunc("/submit", authRequired(submitHandler))
	// http.HandleFunc("/edit/", authRequired(editHandler))
	// http.HandleFunc("/delete/", deleteHandler)
	// http.HandleFunc("/register", registerUserHandler)
	// http.HandleFunc("/login", loginUserHandler)
	// http.HandleFunc("/logout", logoutUserHandler)
	// http.HandleFunc("/oops", status500Handler)

	// http.ListenAndServe(":8000", nil)
}

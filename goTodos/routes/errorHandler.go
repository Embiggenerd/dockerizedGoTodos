package routes

import (
	"fmt"
	"net/http"
)

func errorHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("status500Handler Invoked")

	err := tmplts.ExecuteTemplate(w, "index.html", templData{State: "error", Header: "There was an error", Styles: cacheBustedCss, TodoId: "", Todos: nil, User: nil})
	if err != nil {
		fmt.Println(err)
	}
}

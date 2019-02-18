package routes

import "net/http"

func status500Handler(w http.ResponseWriter, r *http.Request) {
	tmplts.ExecuteTemplate(w, "500.html", nil)
}

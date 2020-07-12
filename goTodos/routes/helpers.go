package routes

import (
	"goTodos/utils"
	"net/http"
)

// RespondWithError sends an appropriate error message and writes to the template
func respondWithError(
	w http.ResponseWriter,
	templateState string,
	err *utils.HTTPError,
) {

	w.WriteHeader(http.StatusUnprocessableEntity)
	tmplts.ExecuteTemplate(w, "index.html", templData{
		State:  templateState,
		Header: "Register with an email and password",
		TodoId: "",
		Todos:  nil,
		User:   nil,
		Error:  err,
	})
}

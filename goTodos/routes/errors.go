package routes

import (
	"goTodos/utils"
	"net/http"
)

// Error handles status codes
type Error struct {
	Code int
	Msg  string
}

func (e Error) Error() string {
	return e.Msg
}

// RespondWithError sends an appropriate error message and writes to the template
func (e Error) RespondWithError(w http.ResponseWriter, code int, msg, templateState string) {

	e.Code = code
	e.Msg = msg

	w.WriteHeader(http.StatusUnprocessableEntity)
	tmplts.ExecuteTemplate(w, "index.html", templData{
		State:  templateState,
		Header: "Register with an email and password",
		TodoId: "",
		Todos:  nil,
		User:   nil,
		Error: &utils.HTTPError{
			Code: e.Code,
			Name: "",
			Msg:  e.Msg,
		},
	})
}

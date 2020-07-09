package utils

import (
	"fmt"
	"net/http"
)

// HTTPError handles status codes
type HTTPError struct {
	Code int
	Name string
	Msg  string
}

func (e HTTPError) Error() string {
	return e.Msg
}

// InternalServerError handles errors with code 500
func InternalServerError(w http.ResponseWriter, r *http.Request) {

	fmt.Println("InternServerError invoked")
}

// UnauthorizedUserError handles errors with code 401
func UnauthorizedUserError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte("Need authorization"))
}

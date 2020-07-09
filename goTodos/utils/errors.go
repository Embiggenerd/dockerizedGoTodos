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
	// w.WriteHeader(http.StatusInternalServerError)
	fmt.Println("InternServerError invoked")
	http.Redirect(w, r, "/oops", http.StatusInternalServerError)
}

// UnauthorizedUserError handles errors with code 401
func UnauthorizedUserError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte("Need authorization"))
}

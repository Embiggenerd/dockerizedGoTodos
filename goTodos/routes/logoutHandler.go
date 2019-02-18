package routes

import "net/http"

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

package auth

import (
	"database/sql"
	"fmt"
	handlerfuncitons "forum/handlers"
	"net/http"
	"forum/database"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		fmt.Println("error executing:", err)
	} else {
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		fmt.Println("error opening database:", err)
	}
	defer db.Close()
	username, err := database.GetUsernameFromToken(db, cookie.Value)
	if err != nil {
		handlerfuncitons.InternalServerError(w, r)
		return
	}
	_, err = db.Exec("DELETE FROM sessions WHERE username = ?", username)
	if err != nil {
		fmt.Println(err)
	}}
	cookie = &http.Cookie{
		Name:   "session_token",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

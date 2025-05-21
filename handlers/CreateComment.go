package handlerfuncitons

import (
	"database/sql"
	"fmt"
	"forum/database"
	"net/http"
	"strconv"
	"strings"
)

func CreateComment(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		fmt.Println("Error in : ", err)
		InternalServerError(w, r)
		return
	}
	defer db.Close()

	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	username, err := database.GetUsernameFromToken(db, cookie.Value)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	var gender string
	err = db.QueryRow("SELECT gender FROM User WHERE username = ?", username).Scan(&gender)
	if err != nil {
		fmt.Println("error getting the gender of the comment writer")
		InternalServerError(w, r)
	}
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			InternalServerError(w, r)
			return
		}

		postID := r.FormValue("postID")
		comment := r.FormValue("comment")
		comment = strings.ReplaceAll(comment, "\\n", "<br>")
		pid, err := strconv.Atoi(postID)
		if err != nil {
			http.Error(w, "Invalid Post ID", http.StatusBadRequest)
			return
		}

		err = database.SaveComment(db, pid, username, comment, gender)
		if err != nil {
			fmt.Println("Error in CreateComment function:", err)
			InternalServerError(w, r)
			return
		}
		db.Exec(`UPDATE Post SET commentcount = commentcount + 1 WHERE post_id = ?`, pid)
		http.Redirect(w, r, fmt.Sprintf("/comment?postID=%d", pid), http.StatusSeeOther)
		return
	}

	MethodNotAllowed(w, r)
}

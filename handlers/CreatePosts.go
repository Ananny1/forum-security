package handlerfuncitons

import (
	"database/sql"
	"fmt"
	"forum/database"
	"net/http"
	"strings"
)

func Createpostshandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		fmt.Println("Error opening database:", err)
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
	user, err := database.GetUser(username, db)
	if err != nil {
		fmt.Println("Error in getUser func: ", err)
		return
	}
	if r.Method == "POST" {
		// Parse form data
		err := r.ParseForm()
		if err != nil {
			InternalServerError(w, r)
			return
		}
		title := r.FormValue("title")
		content := r.FormValue("content")
		categories := r.Form["category"] // assuming "category" form fields contain category IDs
		if !validateCategories(categories) {
			BadRequest(w, r)
			return
		}
		// Replace newlines with newline
		content = strings.ReplaceAll(content, "\\n", "<br>")
		err = database.CreatePost(db, user.Username, title, content, categories, user.Gender)
		if err != nil {
			fmt.Println("Error in CreatePost func: ", err)
			InternalServerError(w, r)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

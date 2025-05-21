package handlerfuncitons

import (
	"database/sql"
	"fmt"
	"forum/database"
	"html/template"
	"net/http"
)

func Post(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/post" { // Check if the requested path is the root URL
		filename := "Templates/post.html"
		t, err := template.ParseFiles(filename) // Parse the index.html template file
		if err != nil {
			fmt.Println("error parsing:", err)
			return
		}
		cookie, err := r.Cookie("session_token")
		if err != nil {
			err = t.Execute(w, nil)
			if err != nil {
				fmt.Println("error executing:", err)
				return
			}
			return
		}
		db, err := sql.Open("sqlite3", "./forum.db")
		if err != nil {
			fmt.Println("error opening database:", err)
		}
		defer db.Close()
		username := cookie.Value
		user, err := database.GetUser(username, db)
		if err != nil {
			fmt.Println("Error getting the user: " , err)
			return
		}
		err = t.ExecuteTemplate(w, "post.html", user) // Execute the template and write the output to the response writer
		if err != nil {
			fmt.Println("error executing:", err)
			return
		}
	} else { // Handle 404 error for other URLs
		NotFoundHandler(w, r)
	}
}

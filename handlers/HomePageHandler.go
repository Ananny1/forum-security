package handlerfuncitons

import (
	"database/sql"
	"fmt"
	"forum/database"
	"forum/structs"
	"html/template"
	"net/http"
	"strconv"
)

func Homepage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" { // Check if the requested path is the root URL
		tmpl, err := template.ParseFiles("templates/base.html", "templates/home.html")
		if err != nil {
			fmt.Println("Error in ParseFiles in HomePageHandler: ", err)
			InternalServerError(w, r)
			return
		}

		if err != nil {
			fmt.Println("error parsing:", err)
			return
		}
		// open db to getPosts
		db, err := sql.Open("sqlite3", "./forum.db")
		if err != nil {
			fmt.Println("Error opening database:", err)
		}
		defer db.Close()

		var Data structs.HomepageData

		// getPosts to display them
		Data.Posts, err = database.GetPosts(db)
		if err != nil {
			fmt.Println("Error in GetPosts func: ", err)
		}

		// initially assume user is guest
		Data.CurrentUser.Username = "guest"
		Data.CurrentUser.Gender = "nil"

		// here we check if cookie already exists
		cookie, err := r.Cookie("session_token")
		// if there is no cookie then assume user is user guest
		if err != nil {
			err = nil // to make cookie error nil and notice the execute error
			err = tmpl.ExecuteTemplate(w, "base.html", Data)
			if err != nil {
				fmt.Println("Error in ExecuteTemplate in HomePageHandler: ", err)
				InternalServerError(w, r)
				return
			}
		} else {
			// if cookie exists then do the following
			username, err := database.GetUsernameFromToken(db, cookie.Value)
			// change username to current userName
			if err != nil || username == "" {
				err = nil
				err = tmpl.ExecuteTemplate(w, "base.html", Data)
				if err != nil {
					fmt.Println("Error in ExecuteTemplate in HomePageHandler: ", err)
					InternalServerError(w, r)
					return
				}
				return
			} else {
				user, err := database.GetUser(username, db)
				if err != nil {
					fmt.Println("Error in getUser func: ", err)
					return
				}
				userids, err := database.GetUserID(db, username)
				if err != nil {
					fmt.Println("Error in GetUserID func: ", err)
					return
				}
				userid, _ := strconv.Atoi(userids)
				Data.Posts, err = database.DisLikedpostsdis(db, userid, Data.Posts)
				if err != nil {
					fmt.Println("Error in GetPosts func: ", err)
				}
				Data.Posts, err = database.Likedpostsdis(db, userid, Data.Posts)
				if err != nil {
					fmt.Println("Error in GetPosts func: ", err)
				}
				Data.CurrentUser = user
				err = tmpl.ExecuteTemplate(w, "base.html", Data)
				if err != nil {
					fmt.Println("Error in ExecuteTemplate in HomePageHandler: ", err)
					InternalServerError(w, r)
					return
				}
			}
		}
	} else { // Handle 404 error for other URLs
		NotFoundHandler(w, r)
	}
}

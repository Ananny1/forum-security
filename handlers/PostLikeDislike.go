package handlerfuncitons

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"forum/database"
	"forum/structs"

	//"html/template"
	"log"
	"net/http"
)

func LikesHandler(w http.ResponseWriter, r *http.Request) {

	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			InternalServerError(w, r)
			return
		}

		postID := r.FormValue("postID")

		var user structs.CurrentUser

		// here we check if cookie already exists
		cookie, err := r.Cookie("session_token")
		// if cookie doesn't exist then redirect user to login page
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusUnauthorized)
			return
		} else {
			// get username from the cookie
			user.Username, err = database.GetUsernameFromToken(db, cookie.Value)
			if err != nil {
				if err == sql.ErrNoRows {
					http.Redirect(w, r, "/login", http.StatusUnauthorized) // No user found for the session token
					return
				}
				fmt.Println("Error in getUsernameFromToken function: ", err)
				InternalServerError(w, r)
				return
			}
		}

		// get the user ID using the username
		userID, err := database.GetUserID(db, user.Username)
		if err != nil {
			fmt.Println("Error in GetUserID function: ", err)
			InternalServerError(w, r)
			return
		}

		// Check if the user has already liked the post
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM LikePost WHERE post_id = ? AND user_id = ?", postID, userID).Scan(&count)
		if err != nil {
			InternalServerError(w, r)
			return
		}

		if count > 0 {
			// ADD FUNCTION TO UNLIKE
			// Check if the user has liked the post, if so, remove the like
			_, err = db.Exec("DELETE FROM LikePost WHERE post_id = ? AND user_id = ?", postID, userID)
			if err != nil {
				InternalServerError(w, r)
				return
			}
			_, err = db.Exec("UPDATE Post SET likecount = likecount - 1 WHERE post_id = ?", postID)
			if err != nil {
				InternalServerError(w, r)
				return
			}

			var likeCount int
			err = db.QueryRow("SELECT likecount FROM Post WHERE post_id = ?", postID).Scan(&likeCount)
			if err != nil {
				InternalServerError(w, r)
				return
			}

			var dislikeCount int
			err = db.QueryRow("SELECT dislikecount FROM Post WHERE post_id = ?", postID).Scan(&dislikeCount)
			if err != nil {
				InternalServerError(w, r)
				return
			}

			response := map[string]interface{}{
				"likedByUser":  false,
				"success":      true,
				"dislikeCount": dislikeCount,
				"likeCount":    likeCount,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}

		// Check if the user has disliked the post, if so, remove the dislike
		err = db.QueryRow("SELECT COUNT(*) FROM DisLikePost WHERE post_id = ? AND user_id = ?", postID, userID).Scan(&count)
		if err != nil {
			InternalServerError(w, r)
			return
		}

		if count > 0 {
			_, err = db.Exec("DELETE FROM DisLikePost WHERE post_id = ? AND user_id = ?", postID, userID)
			if err != nil {
				InternalServerError(w, r)
				return
			}
			_, err = db.Exec("UPDATE Post SET dislikecount = dislikecount - 1 WHERE post_id = ?", postID)
			if err != nil {
				InternalServerError(w, r)
				return
			}
		}

		_, err = db.Exec("INSERT INTO LikePost (post_id, user_id) VALUES (?, ?)", postID, userID)
		if err != nil {
			InternalServerError(w, r)
			return
		}

		_, err = db.Exec("UPDATE Post SET likecount = likecount + 1 WHERE post_id = ?", postID)
		if err != nil {
			InternalServerError(w, r)
			return
		}

		var likeCount int
		err = db.QueryRow("SELECT likecount FROM Post WHERE post_id = ?", postID).Scan(&likeCount)
		if err != nil {
			InternalServerError(w, r)
			return
		}

		var dislikeCount int
		err = db.QueryRow("SELECT dislikecount FROM Post WHERE post_id = ?", postID).Scan(&dislikeCount)
		if err != nil {
			InternalServerError(w, r)
			return
		}

		response := map[string]interface{}{
			"success":        true,
			"likeCount":      likeCount,
			"dislikeCount":   dislikeCount,
			"likedByUser":    true,
			"dislikedbyuser": false,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		MethodNotAllowed(w, r)
	}
}

func DisLikesHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			InternalServerError(w, r)
			return
		}

		postID := r.FormValue("postID")

		var user structs.CurrentUser

		// here we check if cookie already exists
		cookie, err := r.Cookie("session_token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusUnauthorized)
		} else {
			// get username from the cookie
			user.Username, err = database.GetUsernameFromToken(db, cookie.Value)
			if err != nil {
				fmt.Println("Error in GetUsernameFromToken in DislikeHandler: ", err)
				InternalServerError(w, r)
				return
			}
		}

		// get the user ID using the username
		userID, err := database.GetUserID(db, user.Username)
		if err != nil {
			fmt.Println("Error in GetUserID in DislikeHandler: ", err)
			InternalServerError(w, r)
			return
		}

		// Check if the user has already disliked the post
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM DisLikePost WHERE post_id = ? AND user_id = ?", postID, userID).Scan(&count)
		if err != nil {
			InternalServerError(w, r)
			return
		}

		if count > 0 {
			// ADD FUNCTION TO UNDISLIKE
			// Check if the user has disliked the post, if so, remove the dislike
			_, err = db.Exec("DELETE FROM DisLikePost WHERE post_id = ? AND user_id = ?", postID, userID)
			if err != nil {
				InternalServerError(w, r)
				return
			}
			_, err = db.Exec("UPDATE Post SET dislikecount = dislikecount - 1 WHERE post_id = ?", postID)
			if err != nil {
				InternalServerError(w, r)
				return
			}

			var likeCount int
			err = db.QueryRow("SELECT likecount FROM Post WHERE post_id = ?", postID).Scan(&likeCount)
			if err != nil {
				InternalServerError(w, r)
				return
			}

			var dislikeCount int
			err = db.QueryRow("SELECT dislikecount FROM Post WHERE post_id = ?", postID).Scan(&dislikeCount)
			if err != nil {
				InternalServerError(w, r)
				return
			}

			// User has already disliked the post, do nothing or provide feedback
			//w.WriteHeader(http.StatusBadRequest)
			response := map[string]interface{}{
				"dislikedbyuser": false,
				"success":        true,
				"dislikeCount":   dislikeCount,
				"likeCount":      likeCount,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}

		// Check if the user has liked the post, if so, remove the like
		err = db.QueryRow("SELECT COUNT(*) FROM LikePost WHERE post_id = ? AND user_id = ?", postID, userID).Scan(&count)
		if err != nil {
			InternalServerError(w, r)
			return
		}

		if count > 0 {
			_, err = db.Exec("DELETE FROM LikePost WHERE post_id = ? AND user_id = ?", postID, userID)
			if err != nil {
				InternalServerError(w, r)
				return
			}
			_, err = db.Exec("UPDATE Post SET likecount = likecount - 1 WHERE post_id = ?", postID)
			if err != nil {
				InternalServerError(w, r)
				return
			}
		}

		_, err = db.Exec("INSERT INTO DisLikePost (post_id, user_id) VALUES (?, ?)", postID, userID)
		if err != nil {
			InternalServerError(w, r)
			return
		}

		_, err = db.Exec("UPDATE Post SET dislikecount = dislikecount + 1 WHERE post_id = ?", postID)
		if err != nil {
			InternalServerError(w, r)
			return
		}

		var dislikeCount int
		err = db.QueryRow("SELECT dislikecount FROM Post WHERE post_id = ?", postID).Scan(&dislikeCount)
		if err != nil {
			InternalServerError(w, r)
			return
		}

		var likeCount int
		err = db.QueryRow("SELECT likecount FROM Post WHERE post_id = ?", postID).Scan(&likeCount)
		if err != nil {
			InternalServerError(w, r)
			return
		}

		response := map[string]interface{}{
			"success":        true,
			"dislikeCount":   dislikeCount,
			"likeCount":      likeCount,
			"likedByUser":    false,
			"dislikedbyuser": true,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		/*tmpl, err := template.ParseFiles("Templates/register.html")
		if err != nil {
			InternalServerError(w, r)
			return
		}
		err = tmpl.Execute(w, nil)
		if err != nil {
			InternalServerError(w, r)
			return
		}*/
		MethodNotAllowed(w, r)
	}
}

package handlerfuncitons

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"forum/database"
	"forum/structs"
	"log"
	"net/http"
)

func CommentLikeHandler(w http.ResponseWriter, r *http.Request) {
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

		commentID := r.FormValue("commentID")

		var user structs.CurrentUser
		cookie, err := r.Cookie("session_token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusUnauthorized)
			return
		}

		user.Username, err = database.GetUsernameFromToken(db, cookie.Value)
		if err != nil {
			fmt.Println("Error in getUsernameFromToken function: ", err)
			InternalServerError(w, r)
			return
		}

		userID, err := database.GetUserID(db, user.Username)
		if err != nil {
			fmt.Println("Error in GetUserID function: ", err)
			InternalServerError(w, r)
			return
		}

		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM LikeComment WHERE comment_id = ? AND user_id = ?", commentID, userID).Scan(&count)
		if err != nil {
			InternalServerError(w, r)
			return
		}

		// here delete comment like if it was already liked
		if count > 0 {

			_, err = db.Exec("DELETE FROM LikeComment WHERE comment_id = ? AND user_id = ?", commentID, userID)
			if err != nil {
				InternalServerError(w, r)
				return
			}
			_, err = db.Exec("UPDATE Comment SET likecount = likecount - 1 WHERE comment_id = ?", commentID)
			if err != nil {
				InternalServerError(w, r)
				return
			}

			var dislikeCount int
			err = db.QueryRow("SELECT dislikecount FROM Comment WHERE comment_id = ?", commentID).Scan(&dislikeCount)
			if err != nil {
				InternalServerError(w, r)
				return
			}

			var likeCount int
			err = db.QueryRow("SELECT likecount FROM Comment WHERE comment_id = ?", commentID).Scan(&likeCount)
			if err != nil {
				InternalServerError(w, r)
				return
			}

			response := map[string]interface{}{
				"likedByUser":  false,
				"dislikeCount": dislikeCount,
				"likeCount":    likeCount,
				"success":      true,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}

		err = db.QueryRow("SELECT COUNT(*) FROM DisLikeComment WHERE comment_id = ? AND user_id = ?", commentID, userID).Scan(&count)
		if err != nil {
			InternalServerError(w, r)
			return
		}

		if count > 0 {
			_, err = db.Exec("DELETE FROM DisLikeComment WHERE comment_id = ? AND user_id = ?", commentID, userID)
			if err != nil {
				InternalServerError(w, r)
				return
			}
			_, err = db.Exec("UPDATE Comment SET dislikecount = dislikecount - 1 WHERE comment_id = ?", commentID)
			if err != nil {
				InternalServerError(w, r)
				return
			}
		}

		_, err = db.Exec("INSERT INTO LikeComment (comment_id, user_id) VALUES (?, ?)", commentID, userID)
		if err != nil {
			InternalServerError(w, r)
			return
		}

		_, err = db.Exec("UPDATE Comment SET likecount = likecount + 1 WHERE comment_id = ?", commentID)
		if err != nil {
			InternalServerError(w, r)
			return
		}

		var likeCount int
		err = db.QueryRow("SELECT likecount FROM Comment WHERE comment_id = ?", commentID).Scan(&likeCount)
		if err != nil {
			InternalServerError(w, r)
			return
		}

		var dislikeCount int
		err = db.QueryRow("SELECT dislikecount FROM Comment WHERE comment_id = ?", commentID).Scan(&dislikeCount)
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

func CommentDisLikeHandler(w http.ResponseWriter, r *http.Request) {
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

		commentID := r.FormValue("commentID")

		var user structs.CurrentUser
		cookie, err := r.Cookie("session_token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusUnauthorized)
			return
		}

		user.Username, err = database.GetUsernameFromToken(db, cookie.Value)
		if err != nil {
			fmt.Println("Error in getUsernameFromToken function: ", err)
			InternalServerError(w, r)
			return
		}

		userID, err := database.GetUserID(db, user.Username)
		if err != nil {
			fmt.Println("Error in GetUserID function: ", err)
			InternalServerError(w, r)
			return
		}

		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM DisLikeComment WHERE comment_id = ? AND user_id = ?", commentID, userID).Scan(&count)
		if err != nil {
			InternalServerError(w, r)
			return
		}

		if count > 0 {

			_, err = db.Exec("DELETE FROM DisLikeComment WHERE comment_id = ? AND user_id = ?", commentID, userID)
			if err != nil {
				InternalServerError(w, r)
				return
			}
			_, err = db.Exec("UPDATE Comment SET dislikecount = dislikecount - 1 WHERE comment_id = ?", commentID)
			if err != nil {
				InternalServerError(w, r)
				return
			}

			var dislikeCount int
			err = db.QueryRow("SELECT dislikecount FROM Comment WHERE comment_id = ?", commentID).Scan(&dislikeCount)
			if err != nil {
				InternalServerError(w, r)
				return
			}

			var likeCount int
			err = db.QueryRow("SELECT likecount FROM Comment WHERE comment_id = ?", commentID).Scan(&likeCount)
			if err != nil {
				InternalServerError(w, r)
				return
			}

			response := map[string]interface{}{
				"success":        true,
				"dislikedbyuser": false,
				"dislikeCount":   dislikeCount,
				"likeCount":      likeCount,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}

		err = db.QueryRow("SELECT COUNT(*) FROM LikeComment WHERE comment_id = ? AND user_id = ?", commentID, userID).Scan(&count)
		if err != nil {
			InternalServerError(w, r)
			return
		}

		if count > 0 {
			_, err = db.Exec("DELETE FROM LikeComment WHERE comment_id = ? AND user_id = ?", commentID, userID)
			if err != nil {
				InternalServerError(w, r)
				return
			}
			_, err = db.Exec("UPDATE Comment SET likecount = likecount - 1 WHERE comment_id = ?", commentID)
			if err != nil {
				InternalServerError(w, r)
				return
			}
		}

		_, err = db.Exec("INSERT INTO DisLikeComment (comment_id, user_id) VALUES (?, ?)", commentID, userID)
		if err != nil {
			InternalServerError(w, r)
			return
		}

		_, err = db.Exec("UPDATE Comment SET dislikecount = dislikecount + 1 WHERE comment_id = ?", commentID)
		if err != nil {
			InternalServerError(w, r)
			return
		}

		var dislikeCount int
		err = db.QueryRow("SELECT dislikecount FROM Comment WHERE comment_id = ?", commentID).Scan(&dislikeCount)
		if err != nil {
			InternalServerError(w, r)
			return
		}

		var likeCount int
		err = db.QueryRow("SELECT likecount FROM Comment WHERE comment_id = ?", commentID).Scan(&likeCount)
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
		MethodNotAllowed(w, r)
	}
}

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

func Comment(w http.ResponseWriter, r *http.Request) {
	postIDStr := r.URL.Query().Get("postID")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		InternalServerError(w, r)
		return
	}
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		fmt.Println("Error opening database:", err)
	}
	defer db.Close()

	var commentsData structs.CommentPost

	commentsData.Post, err = database.GetPost(db, postID)
	if err != nil {
		fmt.Println("Error getting the post: ", err)
		BadRequest(w, r)
		return
	}

	commentsData.Comments, err = database.GetCommentsByPostID(db, postID)

	if err != nil {
		fmt.Println("Error getting the post: ", err)
		InternalServerError(w, r)
	}
	filename := "Templates/comment.html"
	t, err := template.ParseFiles(filename) // Parse the comment.html template file
	if err != nil {
		fmt.Println("error parsing:", err)
		return
	}
	commentsData.CurrentUser.Username = "guest"
	commentsData.CurrentUser.Gender = "nil"
	cookie, err := r.Cookie("session_token")
	// if there is no cookie then assume user is user guest
	if err != nil {
		err = nil // to make cookie error nil and notice the execute error
		err = t.ExecuteTemplate(w, "comment.html", commentsData)
		if err != nil {
			fmt.Println("error executing:", err)
		}
	} else {
		// if cookie exists then do the following
		username, err := database.GetUsernameFromToken(db, cookie.Value)
		// change username to current userName
		if err != nil || username == "" {
			err = nil
			err = t.ExecuteTemplate(w, "comment.html", commentsData)
			if err != nil {
				fmt.Println(err)
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
			commentsData.Comments, err = database.DisLikedcommsdis(db, userid, commentsData.Comments)
			if err != nil {
				fmt.Println("Error in GetPosts func: ", err)
			}
			commentsData.Comments, err = database.Likedcommsdis(db, userid, commentsData.Comments)
			if err != nil {
				fmt.Println("Error in GetPosts func: ", err)
			}
			commentsData.Post, err = OneDisLikedpostsdis(db, userid, commentsData.Post)
			if err != nil {
				fmt.Println("Error getting the post: ", err)
				BadRequest(w, r)
				return
			}
			commentsData.Post, err = OneLikedpostsdis(db, userid, commentsData.Post)
			if err != nil {
				fmt.Println("Error getting the post: ", err)
				BadRequest(w, r)
				return
			}
			commentsData.CurrentUser = user

			err = t.ExecuteTemplate(w, "comment.html", commentsData)
			// Execute the template and write the output to the response writer
			if err != nil {
				fmt.Println("error executing:", err)
				return
			}
		}
	}

}

func OneDisLikedpostsdis(db *sql.DB, userid int, post structs.Post) (structs.Post, error) {
	rows, err := db.Query(`
        SELECT post_id
        FROM DisLikePost
        WHERE user_id = ?
    `, userid)
	if err != nil {
		return structs.Post{}, fmt.Errorf("error querying posts: %w", err)
	}
	defer rows.Close()

	var dislikedPostID int
	for rows.Next() {
		if err := rows.Scan(&dislikedPostID); err != nil {
			return structs.Post{}, fmt.Errorf("error scanning post ID: %w", err)
		}
		if dislikedPostID == post.ID {
			post.Dislikedbyuser = true
			break
		}
	}

	return post, nil
}

func OneLikedpostsdis(db *sql.DB, userid int, post structs.Post) (structs.Post, error) {
	rows, err := db.Query(`
        SELECT post_id
        FROM LikePost
        WHERE user_id = ?
    `, userid)
	if err != nil {
		return structs.Post{}, fmt.Errorf("error querying posts: %w", err)
	}
	defer rows.Close()

	var likedPostID int
	for rows.Next() {
		if err := rows.Scan(&likedPostID); err != nil {
			return structs.Post{}, fmt.Errorf("error scanning post ID: %w", err)
		}
		if likedPostID == post.ID {
			post.Likedbyuser = true
			break
		}
	}

	return post, nil
}

package handlerfuncitons

import (
	"database/sql"
	"fmt"
	"forum/database"
	"forum/structs"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

func Profilehandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/profile" { // Check if the requested path is the profile URL
		filename := "Templates/ProfilePage.html"
		t, err := template.ParseFiles(filename)
		if err != nil {
			fmt.Println("Error parsing template:", err)
			InternalServerError(w, r)
			return
		}

		// Open DB to get posts
		db, err := sql.Open("sqlite3", "./forum.db")
		if err != nil {
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

		var ProfileData structs.Profiledata
		if r.Method == http.MethodGet { // Use http.MethodGet instead of "Get"
			err := r.ParseForm()
			if err != nil {
				InternalServerError(w, r)
				return
			}
			types := r.FormValue("type")
			if types == "Liked-Post" {
				ProfileData.RequestedPosts, err = GetLikedPosts(db, username)
				userids, err := database.GetUserID(db, username)
				if err != nil {
					fmt.Println("Error in GetUserID func: ", err)
					return
				}
				userid, _ := strconv.Atoi(userids)
				ProfileData.RequestedPosts, err = database.DisLikedpostsdis(db, userid, ProfileData.RequestedPosts)
				if err != nil {
					fmt.Println("Error in GetPosts func: ", err)
				}
				ProfileData.RequestedPosts, err = database.Likedpostsdis(db, userid, ProfileData.RequestedPosts)
				if err != nil {
					fmt.Println("Error in GetPosts func: ", err)
				}
				if err != nil {
					fmt.Println("Error in GetLikedPosts func:", err)
					InternalServerError(w, r)
					return
				}
			} else if types == "Created-Post" {
				ProfileData.RequestedPosts, err = GetCreatedPosts(db, username)
				userids, err := database.GetUserID(db, username)
				if err != nil {
					fmt.Println("Error in GetUserID func: ", err)
					return
				}
				userid, _ := strconv.Atoi(userids)
				ProfileData.RequestedPosts, err = database.DisLikedpostsdis(db, userid, ProfileData.RequestedPosts)
				if err != nil {
					fmt.Println("Error in GetPosts func: ", err)
				}
				ProfileData.RequestedPosts, err = database.Likedpostsdis(db, userid, ProfileData.RequestedPosts)
				if err != nil {
					fmt.Println("Error in GetPosts func: ", err)
				}
				if err != nil {
					fmt.Println("Error in GetCreatedPosts func:", err)
					InternalServerError(w, r)
					return
				}
			} else if types != "" {
				BadRequest(w, r)
				return
			}
			ProfileData.Username = username
			err = db.QueryRow("SELECT gender, email,pfplink FROM User WHERE username = ?", username).Scan(&ProfileData.Gender, &ProfileData.Email,&ProfileData.ProfilePicture)
			if err != nil {
				fmt.Println("Error getting gender and email:", err)
				InternalServerError(w, r)
				return
			}
			err = t.ExecuteTemplate(w, "ProfilePage.html", ProfileData)
			if err != nil {
				fmt.Println("Error executing template:", err)
				InternalServerError(w, r)
				return
			}
		}
	}
}

func GetLikedPosts(db *sql.DB, username string) ([]structs.Post, error) {
	// Step 1: Retrieve user_id based on username from User table
	var userID int
	err := db.QueryRow(`
        SELECT user_id FROM User WHERE username = ?
    `, username).Scan(&userID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving user_id: %w", err)
	}

	// Step 2: Retrieve all post_id from LikePost table based on user_id
	rows, err := db.Query(`
        SELECT post_id FROM LikePost WHERE user_id = ?
    `, userID)
	if err != nil {
		return nil, fmt.Errorf("error querying liked posts: %w", err)
	}
	defer rows.Close()

	// Collect post_id in an array
	var postIDs []int
	for rows.Next() {
		var postID int
		if err := rows.Scan(&postID); err != nil {
			return nil, fmt.Errorf("error scanning post_id: %w", err)
		}
		postIDs = append(postIDs, postID)
	}

	if len(postIDs) == 0 {
		return nil, nil // No liked posts found
	}

	// Step 3: Construct the SQL query manually
	query := `
        SELECT post_id, username, title, content, created_at, likecount, dislikecount, gender, commentcount
        FROM Post WHERE post_id IN (`
	placeholders := make([]string, len(postIDs))
	args := make([]interface{}, len(postIDs))
	for i, id := range postIDs {
		placeholders[i] = "?"
		args[i] = id
	}
	query += strings.Join(placeholders, ",") + `)
        ORDER BY created_at`

	// Execute the query with the collected postIDs
	rows, err = db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error querying posts: %w", err)
	}
	defer rows.Close()

	var posts []structs.Post
	for rows.Next() {
		var post structs.Post
		err := rows.Scan(&post.ID, &post.Username, &post.Title, &post.Content, &post.CreatedAt, &post.Like, &post.Dislike, &post.Gender, &post.Comment)
		if err != nil {
			return nil, fmt.Errorf("error scanning post row: %w", err)
		}
		post.CreatedAt = database.ReplaceLettersWithSpaces(post.CreatedAt)
		if post.Gender == "Google" || post.Gender == "Github" {
			err := db.QueryRow("SELECT pfplink FROM User WHERE username = ?", post.Username).Scan(&post.ProfilePicture)
			if err != nil {
				return posts, err
			}
		}

		// Retrieve categories for the post
		categories, err := database.GetCategoriesForPost(db, post.ID)
		if err != nil {
			return nil, fmt.Errorf("error getting categories for post: %w", err)
		}
		post.Category = categories

		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over posts: %w", err)
	}

	return posts, nil
}

func GetCreatedPosts(db *sql.DB, username string) ([]structs.Post, error) {
	rows, err := db.Query(`
        SELECT post_id, username, title, content, created_at, likecount, dislikecount, gender, commentcount
        FROM Post WHERE username = ?
        ORDER BY created_at
    `, username) // Pass the username here
	if err != nil {
		return nil, fmt.Errorf("error querying posts: %w", err)
	}
	defer rows.Close()

	var posts []structs.Post

	for rows.Next() {
		var post structs.Post
		err := rows.Scan(&post.ID, &post.Username, &post.Title, &post.Content, &post.CreatedAt, &post.Like, &post.Dislike, &post.Gender, &post.Comment)
		if err != nil {
			return nil, fmt.Errorf("error scanning post row: %w", err)
		}
		post.CreatedAt = database.ReplaceLettersWithSpaces(post.CreatedAt)
		if post.Gender == "Google" || post.Gender == "Github" {
			err := db.QueryRow("SELECT pfplink FROM User WHERE username = ?", post.Username).Scan(&post.ProfilePicture)
			if err != nil {
				return posts, err
			}
		}
		// Retrieve categories for the post
		categories, err := database.GetCategoriesForPost(db, post.ID)
		if err != nil {
			return nil, fmt.Errorf("error getting categories for post: %w", err)
		}
		post.Category = categories

		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over posts: %w", err)
	}

	return posts, nil
}

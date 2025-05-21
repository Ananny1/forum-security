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

func FilterHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/base.html", "templates/category.html")
	if err != nil {
		fmt.Println("Error in ParseFiles in CategoryHandler: ", err)
		InternalServerError(w, r)
		return
	}

	var FilteredData structs.CategoryHandlerPage

	// Open db to getPosts
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}
	defer db.Close()

	// get the categories from the DB to display them
	FilteredData.CategoriesList, err = database.GetCategories(db)
	if err != nil {
		fmt.Println("Error in GetCategories function in categoryHandler: ", err)
		InternalServerError(w, r)
		return
	}

	// Get posts to display them
	FilteredData.Posts, err = database.GetPosts(db)
	if err != nil {
		fmt.Println("Error in GetPosts func: ", err)
		return
	}

	// Initially assume user is guest
	FilteredData.CurrentUser.Username = "guest"
	FilteredData.CurrentUser.Gender = "nil"

	if r.Method == "POST" {
		// Parse form data
		err := r.ParseForm()
		if err != nil {
			InternalServerError(w, r)
			return
		}
		categories := r.Form["category"]
		if !validateCategories(categories) {
			BadRequest(w, r)
			return
		}
		FilteredData.Posts = filterPosts(FilteredData.Posts, categories) // Corrected line
	}
	// Check if cookie already exists
	cookie, err := r.Cookie("session_token")
	if err != nil {
		// No cookie, assume user is guest
		err = nil // Reset error to avoid interfering with template execution
		err = tmpl.ExecuteTemplate(w, "base.html", FilteredData)
		if err != nil {
			fmt.Println("Error in ExecuteTemplate in CategoryHandler: ", err)
			InternalServerError(w, r)
			return
		}
		return
	}
	// If cookie exists, get the username from the token
	username, err := database.GetUsernameFromToken(db, cookie.Value)
	if err != nil || username == "" {
		// Invalid token or username, assume guest
		err = nil
		err = tmpl.ExecuteTemplate(w, "base.html", FilteredData)
		if err != nil {
			fmt.Println("Error in ExecuteTemplate in CategoryHandler: ", err)
			InternalServerError(w, r)
			return
		}
		return
	}
	// Valid token, get user details
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
	FilteredData.Posts, err = database.DisLikedpostsdis(db, userid, FilteredData.Posts)
	if err != nil {
		fmt.Println("Error in GetPosts func: ", err)
	}
	FilteredData.Posts, err = database.Likedpostsdis(db, userid, FilteredData.Posts)
	if err != nil {
		fmt.Println("Error in GetPosts func: ", err)
	}

	FilteredData.CurrentUser = user
	err = tmpl.ExecuteTemplate(w, "base.html", FilteredData)
	if err != nil {
		fmt.Println("Error in ExecuteTemplate in CategoryHandler: ", err)
		InternalServerError(w, r)
		return
	}
}

func Contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

// Helper function to check if all categories are present in post categories
func containsAllCategories(postCategories []string, categories []string) bool {
	for _, category := range categories {
		if !Contains(postCategories, category) { // Corrected function call
			return false
		}
	}
	return true
}

// Function to filter posts based on categories
func filterPosts(posts []structs.Post, categories []string) []structs.Post {
	var filteredPosts []structs.Post

	// Iterate through each post
	for _, post := range posts {
		// Check if the post categories contain all the required categories
		if containsAllCategories(post.Category, categories) {
			filteredPosts = append(filteredPosts, post)
		}
	}

	return filteredPosts
}

func Category(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/category" { // Check if the requested path is the root URL
		tmpl, err := template.ParseFiles("templates/base.html", "templates/category.html")
		if err != nil {
			fmt.Println("Error in ParseFiles in categoryHandler: ", err)
			InternalServerError(w, r)
			return
		}
		db, err := sql.Open("sqlite3", "./forum.db")
		if err != nil {
			fmt.Println("error opening database:", err)
		}
		defer db.Close()

		var Data structs.CategoryHandlerPage

		// get the categories from the DB to display them
		Data.CategoriesList, err = database.GetCategories(db)
		if err != nil {
			fmt.Println("Error in GetCategories function in categoryHandler: ", err)
			InternalServerError(w, r)
			return
		}

		Data.Posts, err = database.GetPosts(db)
		if err != nil {
			fmt.Println("Error in GetPosts func: ", err)
		}
		Data.CurrentUser.Username = "guest"
		Data.CurrentUser.Gender = "nil"
		cookie, err := r.Cookie("session_token")
		if err != nil {
			err = tmpl.ExecuteTemplate(w, "base.html", Data)
			if err != nil {
				fmt.Println("Error in ExecuteTemplate in categoryHandler: ", err)
				InternalServerError(w, r)
				return
			}
			return
		} else {
			username, err := database.GetUsernameFromToken(db, cookie.Value)
			// change username to current userName
			if err != nil || username == "" {
				err = nil
				err = tmpl.ExecuteTemplate(w, "base.html", Data)
				if err != nil {
					fmt.Println("Error in ExecuteTemplate in categoryHandler: ", err)
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
					fmt.Println("Error in ExecuteTemplate in categoryHandler: ", err)
					InternalServerError(w, r)
					return
				}
			}
		}
	} else { // Handle 404 error for other URLs
		NotFoundHandler(w, r)
	}
}

func isValidCategory(category string) bool {
	validCategories := map[string]struct{}{
		"Sport":     {},
		"Gaming":    {},
		"Art":       {},
		"Education": {},
		"Food":      {},
	}
	_, exists := validCategories[category]
	return exists
}

// validateCategories checks if all elements in the slice are valid categories
func validateCategories(categories []string) bool {
	for _, category := range categories {
		if !isValidCategory(category) {
			return false
		}
	}
	return true
}

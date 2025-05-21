package database

import (
	"database/sql"
	"fmt"
	"forum/structs"
)

func GetPost(db *sql.DB, postid int) (structs.Post, error) {
	var post structs.Post

	// Query to fetch the post by post_id
	query := `
        SELECT post_id, username, title, content, created_at, likecount, dislikecount, gender, commentcount
        FROM Post
        WHERE post_id = ?
    `
	row := db.QueryRow(query, postid)
	err := row.Scan(&post.ID, &post.Username, &post.Title, &post.Content, &post.CreatedAt, &post.Like, &post.Dislike, &post.Gender, &post.Comment)
	if err != nil {
		return post, err
	}
	if post.Gender == "Google" || post.Gender == "Github" {
		err := db.QueryRow("SELECT pfplink FROM User WHERE username = ?", post.Username).Scan(&post.ProfilePicture)
		if err != nil {
			return post, err
		}
	}

	// Example function to replace letters with spaces (you need to define this function)
	post.CreatedAt = ReplaceLettersWithSpaces(post.CreatedAt)

	// Retrieve categories for the post (you need to define this function)
	categories, err := GetCategoriesForPost(db, post.ID)
	if err != nil {
		return post, fmt.Errorf("error getting categories for post: %w", err)
	}
	post.Category = categories

	return post, nil
}

func GetPosts(db *sql.DB) ([]structs.Post, error) {
	rows, err := db.Query(`
        SELECT post_id, username, title, content, created_at, likecount, dislikecount, gender, commentcount
        FROM Post
        ORDER BY created_at
    `)
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
		post.CreatedAt = ReplaceLettersWithSpaces(post.CreatedAt)
		if post.Gender == "Google" || post.Gender == "Github" {
			err := db.QueryRow("SELECT pfplink FROM User WHERE username = ?", post.Username).Scan(&post.ProfilePicture)
			if err != nil {
				return posts, err
			}
		}
		// Retrieve categories for the post
		categories, err := GetCategoriesForPost(db, post.ID)
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

// Helper function to get categories for a given post
func GetCategoriesForPost(db *sql.DB, postID int) ([]string, error) {
	rows, err := db.Query(`
        SELECT category_id
        FROM PostCategory
        WHERE post_id = ?
    `, postID)
	if err != nil {
		return nil, fmt.Errorf("error querying categories: %w", err)
	}
	defer rows.Close()

	var categoryIDs []string
	for rows.Next() {
		var categoryID string
		err := rows.Scan(&categoryID)
		if err != nil {
			return nil, fmt.Errorf("error scanning category row: %w", err)
		}
		categoryIDs = append(categoryIDs, categoryID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over categories: %w", err)
	}

	return categoryIDs, nil
}

func ReplaceLettersWithSpaces(input string) string {
	var result string
	for _, c := range input {
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') {
			result += " "
		} else {
			result += string(c)
		}
	}
	return result
}

func Likedpostsdis(db *sql.DB, userid int, posts []structs.Post) ([]structs.Post, error) {
	rows, err := db.Query(`
        SELECT post_id
        FROM LikePost
        WHERE user_id = ?
    `, userid)
	if err != nil {
		return nil, fmt.Errorf("error querying posts: %w", err)
	}
	defer rows.Close()

	var likedPostIDs []int
	for rows.Next() {
		var postID int
		if err := rows.Scan(&postID); err != nil {
			return nil, fmt.Errorf("error scanning post ID: %w", err)
		}
		likedPostIDs = append(likedPostIDs, postID)
	}

	likedPostMap := make(map[int]bool)
	for _, id := range likedPostIDs {
		likedPostMap[id] = true
	}

	for i, post := range posts {
		if likedPostMap[post.ID] {
			posts[i].Likedbyuser = true
		}
	}

	return posts, nil
}

func DisLikedpostsdis(db *sql.DB, userid int, posts []structs.Post) ([]structs.Post, error) {
	rows, err := db.Query(`
        SELECT post_id
        FROM DisLikePost
        WHERE user_id = ?
    `, userid)
	if err != nil {
		return nil, fmt.Errorf("error querying posts: %w", err)
	}
	defer rows.Close()

	var DislikedPostIDs []int
	for rows.Next() {
		var postID int
		if err := rows.Scan(&postID); err != nil {
			return nil, fmt.Errorf("error scanning post ID: %w", err)
		}
		DislikedPostIDs = append(DislikedPostIDs, postID)
	}

	likedPostMap := make(map[int]bool)
	for _, id := range DislikedPostIDs {
		likedPostMap[id] = true
	}

	for i, post := range posts {
		if likedPostMap[post.ID] {
			posts[i].Dislikedbyuser = true
		}
	}

	return posts, nil
}

func Likedcommsdis(db *sql.DB, userid int, comments []structs.Comment) ([]structs.Comment, error) {
	rows, err := db.Query(`
        SELECT comment_id
        FROM LikeComment
        WHERE user_id = ?
    `, userid)
	if err != nil {
		return nil, fmt.Errorf("error querying posts: %w", err)
	}
	defer rows.Close()

	var likedCommIDs []int
	for rows.Next() {
		var commID int
		if err := rows.Scan(&commID); err != nil {
			return nil, fmt.Errorf("error scanning post ID: %w", err)
		}
		likedCommIDs = append(likedCommIDs, commID)
	}

	likedCommtMap := make(map[int]bool)
	for _, id := range likedCommIDs {
		likedCommtMap[id] = true
	}

	for i, comment := range comments {
		if likedCommtMap[comment.CommentID] {
			comments[i].Likedbyuser = true
		}
	}

	return comments, nil
}

func DisLikedcommsdis(db *sql.DB, userid int, comments []structs.Comment) ([]structs.Comment, error) {
	rows, err := db.Query(`
        SELECT comment_id
        FROM DisLikeComment
        WHERE user_id = ?
    `, userid)
	if err != nil {
		return nil, fmt.Errorf("error querying posts: %w", err)
	}
	defer rows.Close()

	var dislikedCommIDs []int
	for rows.Next() {
		var commID int
		if err := rows.Scan(&commID); err != nil {
			return nil, fmt.Errorf("error scanning post ID: %w", err)
		}
		dislikedCommIDs = append(dislikedCommIDs, commID)
	}

	dislikedCommtMap := make(map[int]bool)
	for _, id := range dislikedCommIDs {
		dislikedCommtMap[id] = true
	}

	for i, comment := range comments {
		if dislikedCommtMap[comment.CommentID] {
			comments[i].Dislikedbyuser = true
		}
	}

	return comments, nil
}

// to get the catigories and display them in the filter

func GetCategories(db *sql.DB) ([]string, error) {
    // Prepare the query
    query := "SELECT name FROM Category"

    // Execute the query
    rows, err := db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    // Create a slice to hold the categories
    var categories []string

    // Loop through the rows
    for rows.Next() {
        var category string
        if err := rows.Scan(&category); err != nil {
            return nil, err
        }
        categories = append(categories, category)
    }

    // Check for errors during the iteration
    if err := rows.Err(); err != nil {
        return nil, err
    }

    return categories, nil
}

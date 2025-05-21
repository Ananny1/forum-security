package database

import (
	"database/sql"
	"forum/structs"
)

func GetCommentsByPostID(db *sql.DB, postID int) ([]structs.Comment, error) {
	// Prepare the SQL query to get comments by post_id
	query := `SELECT comment_id, post_id, usernames, content, created_at, gender, likecount, dislikecount FROM Comment WHERE post_id = ?`
	// Execute the query
	rows, err := db.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Slice to hold the results
	var comments []structs.Comment

	// Iterate over the rows and scan them into Comment structs
	for rows.Next() {
		var c structs.Comment
		if err := rows.Scan(&c.CommentID, &c.PostID, &c.Username, &c.Content, &c.CreatedAt, &c.Gender, &c.Like, &c.Dislike); err != nil {
			return nil, err
		}
		if c.Gender == "Google" || c.Gender == "Github" {
			err := db.QueryRow("SELECT pfplink FROM User WHERE username = ?", c.Username).Scan(&c.ProfilePicture)
			if err != nil {
				return nil, err
			}
		}
		c.CreatedAt = ReplaceLettersWithSpaces(c.CreatedAt)
		comments = append(comments, c)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

/*// func to get comments by postID
func GetCommentsByPostID(db *sql.DB, postID int) ([]structs.Comment, error) {
	rows, err := db.Query(`
		SELECT comment_id, post_id, username, content, created_at
		FROM Comment
		WHERE post_id = ?
		ORDER BY created_at
	`, postID)
	if err != nil {
		return nil, fmt.Errorf("error querying comments: %w", err)
	}
	defer rows.Close()

	var comments []structs.Comment

	for rows.Next() {
		var comment structs.Comment
		err := rows.Scan(&comment.CommentID, &comment.PostID, &comment.Username, &comment.Content, &comment.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning comment row: %w", err)
		}
		comment.CreatedAt = ReplaceLettersWithSpaces(comment.CreatedAt)
		comments = append(comments, comment)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over comments: %w", err)
	}

	return comments, nil
}*/

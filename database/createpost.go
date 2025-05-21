package database
import (
	"database/sql"
	"fmt"
	"time"
	"forum/structs"
)
func CreatePost(db *sql.DB, username, title, content string, categories []string, gender string) error {
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error beginning transaction: %w", err)
	}
	// Insert into the Post table
	stmtPost, err := tx.Prepare(`
		INSERT INTO Post (username, title, content, created_at, gender) 
		VALUES (?, ?, ?, ?, ?)
	`)
	if err != nil {
		tx.Rollback() // undo all changes made within that transaction
		return fmt.Errorf("error preparing statement for Post: %w", err)
	}
	defer stmtPost.Close()
	// Executing Post Insertion (result var holds information about the execution)
	result, err := stmtPost.Exec(username, title, content, currentTime, gender)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error executing statement for Post: %w", err)
	}
	// Get the ID of the post we just inserted from the result of execution (if there is error print it)
	postID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error getting last insert id: %w", err)
	}
	// Insert the categories into the PostCategory table
	stmtCategory, err := tx.Prepare(`
		INSERT INTO PostCategory (post_id, category_id)
		VALUES (?, ?)
	`)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error preparing statement for PostCategory: %w", err)
	}
	defer stmtCategory.Close()
	//  iterates through the categories
	for _, categoryID := range categories {
		_, err = stmtCategory.Exec(postID, categoryID) // prepared statement is executed for each category
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error executing statement for PostCategory: %w", err)
		}
	}
	// Commit the transaction (making all changes permanent)
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}
	// if there is no error then return nil
	return nil
}
/*	Here we used transactions because it is important to rollback if an issue happend	*/
func SavePost(db *sql.DB, post structs.Post) error {
    query := `INSERT INTO Post (user_id, title, content) VALUES (?, ?, ?)`
    _, err := db.Exec(query, post.Username, post.Title, post.Content)
    return err
}
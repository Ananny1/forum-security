package database

import (
    "database/sql"
    "fmt"
    "time"
)


func SaveComment(db *sql.DB, postID int, username string, content string, gender string) error {
    currentTime := time.Now().Format("2006-01-02 15:04:05")
    s, err := db.Prepare(`
        INSERT INTO Comment (post_id, usernames, content, created_at, gender)
        VALUES (?, ?, ?, ?, ?)
    `)
    if err != nil {
        return fmt.Errorf("error preparing statement: %w", err)
    }
    defer s.Close()

    _, err = s.Exec(postID, username, content, currentTime, gender)
    if err != nil {
        return fmt.Errorf("error executing statement: %w", err)
    }
    return nil
}

/*  Note: Here we didn't use transactions (db.Begin() and tx.Commit() or tx.Rollback()) 
    since there are no multiple database operations need to be treated as a single unit. 
    Here we only perform a single INSERT statement, so if this operation fails, the comment 
    won't be saved anyway.
*/
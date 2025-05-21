package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
)

func Database() {
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		fmt.Println("error opening database:", err)
	}
	defer db.Close()
	// Create tables
	CreateTables(db)
	insertCategories(db)
}

// This function creates the tables if they don't exists
func CreateTables(db *sql.DB) {
	UserTable := `
    CREATE TABLE IF NOT EXISTS User (
        user_id INTEGER PRIMARY KEY AUTOINCREMENT,
        email TEXT UNIQUE NOT NULL,
        username TEXT UNIQUE NOT NULL,
        password TEXT NOT NULL,
		gender TEXT NOT NULL,
		pfplink TEXT
    );`

	PostTable := `
	CREATE TABLE IF NOT EXISTS Post (
		post_id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		gender TEXT NOT NULL,
		likecount INTEGER DEFAULT 0,
		dislikecount INTEGER DEFAULT 0,
		commentcount INTEGER DEFAULT 0,
		image_url TEXT, -- New column for post image
		FOREIGN KEY (username) REFERENCES User(username)
	);
`

	CommentTable := `
    CREATE TABLE IF NOT EXISTS Comment (
        comment_id INTEGER PRIMARY KEY AUTOINCREMENT,
        post_id INTEGER NOT NULL,
		usernames TEXT NOT NULL,
        content TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		gender TEXT NOT NULL,
		likecount INTEGER DEFAULT 0,
		dislikecount INTEGER DEFAULT 0,
        FOREIGN KEY (post_id) REFERENCES Post(post_id),
		FOREIGN KEY (usernames) REFERENCES User(username)
    );`

	CategoryTable := `
    CREATE TABLE IF NOT EXISTS Category (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT UNIQUE NOT NULL
    );`

	PostCategoryTable := `
    CREATE TABLE IF NOT EXISTS PostCategory (
        post_id INTEGER NOT NULL,
        category_id INTEGER NOT NULL,
        FOREIGN KEY (post_id) REFERENCES Post(post_id),
        FOREIGN KEY (category_id) REFERENCES Category(id),
        PRIMARY KEY (post_id, category_id)
    );`

	LikePostTable := `
    CREATE TABLE IF NOT EXISTS LikePost (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER NOT NULL,
        post_id INTEGER NOT NULL,
        FOREIGN KEY (user_id) REFERENCES User(user_id),
        FOREIGN KEY (post_id) REFERENCES Post(post_id)
    );`

	LikeCommentTable := `
    CREATE TABLE IF NOT EXISTS LikeComment (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER NOT NULL,
        comment_id INTEGER NOT NULL,
        FOREIGN KEY (user_id) REFERENCES User(user_id),
        FOREIGN KEY (comment_id) REFERENCES Comment(comment_id)
    );`

	DisLikePostTable := `
    CREATE TABLE IF NOT EXISTS DisLikePost (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER NOT NULL,
        post_id INTEGER NOT NULL,
        FOREIGN KEY (user_id) REFERENCES User(user_id),
        FOREIGN KEY (post_id) REFERENCES Post(post_id)
    );`

	DisLikeCommentTable := `
    CREATE TABLE IF NOT EXISTS DisLikeComment (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER NOT NULL,
        comment_id INTEGER NOT NULL,
        FOREIGN KEY (user_id) REFERENCES User(user_id),
        FOREIGN KEY (comment_id) REFERENCES Comment(comment_id)
    );`

	SessionTable := `
    CREATE TABLE IF NOT EXISTS sessions (
    username VARCHAR(255) NOT NULL,
    token VARCHAR(32) NOT NULL,
    PRIMARY KEY (username)
);`

	// Execute the SQL statements to create tables
	_, err := db.Exec(UserTable)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(PostTable)
	if err != nil {
		log.Fatal("here", err)
	}

	_, err = db.Exec(CommentTable)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(CategoryTable)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(PostCategoryTable)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(LikePostTable)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(LikeCommentTable)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(DisLikePostTable)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(DisLikeCommentTable)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(SessionTable)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Tables created successfully!")
}

// this function adds the categories we chose to the category table
func insertCategories(db *sql.DB) {
	categories := []string{"Sport", "Gaming", "Art", "Education", "Food"}

	for _, category := range categories {
		_, err := db.Exec("INSERT INTO Category (name) VALUES (?)", category)
		if err != nil {
			if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.Code == sqlite3.ErrConstraint {
				// Skip insertion if category already exists
				continue
			} else {
				log.Fatal(err)
			}
		}
	}
}

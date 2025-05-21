package database

import (
	"database/sql"
	"forum/structs"
)

func GetUserID(db *sql.DB, username string) (string, error) {
	var userID string
	err := db.QueryRow("SELECT user_id FROM User WHERE username = ?", username).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil // No user found for the session token
		}
		return "", err
	}
	return userID, nil
}

func GetUser(username string, db *sql.DB) (structs.CurrentUser, error) {
	var user structs.CurrentUser
	err := db.QueryRow("SELECT username, gender, pfplink FROM User WHERE username = ?", username).Scan(&user.Username, &user.Gender, &user.ProfilePicture)
	if err != nil {
		return user, err
	}
	return user, nil
}

func GetUsernameFromToken(db *sql.DB, token string) (string, error) {
	var username string
	err := db.QueryRow("SELECT username FROM sessions WHERE token = ?", token).Scan(&username)
	if err != nil {
		return "", err
	}
	return username, nil
}

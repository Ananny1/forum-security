package auth

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
)

func GenerateSessionToken() (string, error) {
	token := make([]byte, 16)
	if _, err := rand.Read(token); err != nil {
		return "", err
	}
	return hex.EncodeToString(token), nil
}

func SetSessionToken(db *sql.DB, username, token string) error {
	// Invalidate old session
	_, err := db.Exec("DELETE FROM sessions WHERE username = ?", username)
	if err != nil {
		return err
	}

	// Set new session token
	_, err = db.Exec("INSERT INTO sessions (username, token) VALUES (?, ?)", username, token)
	return err
}

func CheckUniqueUsername(username string) (string, error) {
	// open database
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// check if username duplicated
	stmt := "SELECT username FROM User WHERE username = ?"
	row := db.QueryRow(stmt, username)
	var u string
	err = row.Scan(&u)

	// if user is not found in the DataBase just return the name
	if err != nil {
		if err == sql.ErrNoRows {
			return username, nil
		} else {
			fmt.Println("Error: ", err)
			return "", err
		}
	} else {
		// if username already found then generate a unique one for the google user
		isDuplicated := true
		numberCount := 0
		numberString := strconv.Itoa(numberCount)
		NewName := username + numberString

		for isDuplicated {
			query := `SELECT username FROM User WHERE username LIKE ?`
			row := db.QueryRow(query, username)
			var result string
			err := row.Scan(&result)
			if err != nil && err != sql.ErrNoRows {
				numberCount++
			} else {
				isDuplicated = false
			}
		}
		return NewName, nil
	}
}


func isEmailRegisteredInDatabase(email string) (bool, error) {
	// open database
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	email = strings.TrimSpace(email)
	// check if email not empty
	if email == ""{
		return false , errors.New("Empty email")
	}
	// check if email already exists in database
	stmt := "SELECT username FROM User WHERE email = ?"
	row := db.QueryRow(stmt, email)
	var uID string
	err = row.Scan(&uID)

	// if didn't find the email return false
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		} else {
			return false, err
		}
	}

	// if email found return true
	return true, nil
}
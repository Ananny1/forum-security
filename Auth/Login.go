package auth

import (
	"database/sql"
	"fmt"
	"forum/database"
	handlerfuncitons "forum/handlers"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// open database
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	cookie, err := r.Cookie("session_token")
	// if there is no cookie then assume user is user guest
	if err == nil { 
		// if cookie exists then do the following
		_, err := database.GetUsernameFromToken(db, cookie.Value)
		// change username to current userName
		if err == nil {
			http.Redirect(w, r, "/welcome", http.StatusFound)
			return
		}
	}
	if r.Method == http.MethodPost {
		// Parse form data
		err := r.ParseForm()
		if err != nil {
			handlerfuncitons.InternalServerError(w, r)
			return
		}

		username := r.FormValue("username")
		password := r.FormValue("password")

		username = strings.TrimSpace(username)
		password = strings.TrimSpace(password)

		// This is to handle the reddirection from the like and also to validate
		if username == "" || password == "" {
			tmpl, err := template.ParseFiles("Templates/login.html")
			if err != nil {
				fmt.Println(err)
				handlerfuncitons.InternalServerError(w, r)
				return
			}
			tmpl.ExecuteTemplate(w, "login.html", nil)
			return
		}

		var storedHashedPassword string
		err = db.QueryRow("SELECT password FROM User WHERE username = ?", username).Scan(&storedHashedPassword) //query pass and store it
		if err != nil {
			if err == sql.ErrNoRows {
				w.WriteHeader(http.StatusUnauthorized)
				tmpl, err := template.ParseFiles("Templates/login.html")
				if err != nil {
					fmt.Println(err)
					handlerfuncitons.InternalServerError(w, r)
					return
				}
				tmpl.ExecuteTemplate(w, "login.html", "invalidUser")
				return
			} else {
				log.Printf("Database error: %v", err) // Log the detailed database error
				if err != nil {
					handlerfuncitons.InternalServerError(w, r)
					return
				}
			}
			return
		}
		// Compare the password
		err = bcrypt.CompareHashAndPassword([]byte(storedHashedPassword), []byte(password))
		if err != nil {
			// if password isn't correct then user unauthorized
			w.WriteHeader(http.StatusUnauthorized)
			tmpl, err := template.ParseFiles("Templates/login.html")
			if err != nil {
				fmt.Println(err)
			 	handlerfuncitons.InternalServerError(w, r)
				return
			}
			tmpl.ExecuteTemplate(w, "login.html", "invalidUser")
			return
		} else {
			// if username and password correct then add cookie and go to home page
			database.GetUser(username, db)
			sessionToken, err := GenerateSessionToken() // function to generate token
			if err != nil {
				handlerfuncitons.InternalServerError(w, r)
				return
			}
			if err := SetSessionToken(db, username, sessionToken); err != nil {
				handlerfuncitons.InternalServerError(w, r)
				return
			}
			http.SetCookie(w, &http.Cookie{
				Name:     "session_token",
				Value:    sessionToken,
				Path:     "/",
				Expires:  time.Now().Add(24 * time.Hour),
				HttpOnly: true, // cant be accsessd by frontend
				Secure:   true, // will only be sent over https
			})
			// if successful then redirect user to main page
			http.Redirect(w, r, "/welcome", http.StatusFound)
		}
	} else { // Display the form
		tmpl, err := template.ParseFiles("Templates/login.html")
		if err != nil {
			handlerfuncitons.InternalServerError(w, r)
			return
		}
		err = tmpl.Execute(w, nil)
		if err != nil {
			handlerfuncitons.InternalServerError(w, r)
			return
		}
	}
}

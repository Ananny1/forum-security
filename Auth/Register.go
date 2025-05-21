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

	"golang.org/x/crypto/bcrypt"
)

func Registerhandler(w http.ResponseWriter, r *http.Request) {
	// open and connect to the SQLite database
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	cookie, err := r.Cookie("session_token")
	if err == nil {
		_, err := database.GetUsernameFromToken(db, cookie.Value)
		if err == nil {
			http.Redirect(w, r, "/welcome", http.StatusFound)
			return
		}
	}

	// check if method is post (user pressed submit and entered data)
	if r.Method == http.MethodPost {
		// Parse form data
		err := r.ParseForm()
		if err != nil {
			handlerfuncitons.InternalServerError(w, r)
			return
		}

		// get data from form
		email := r.FormValue("email")
		username := r.FormValue("username")
		password := r.FormValue("password")
		gender := r.FormValue("gender")
		pfp:=""
		if gender == "male" {
			pfp = "static/img/male.webp"
		} else if gender == "female" {
			pfp = "static/img/female.png"
		} else {
			handlerfuncitons.BadRequest(w, r)
		}

		// display the registration form
		tmpl, err := template.ParseFiles("Templates/register.html")
		if err != nil {
			handlerfuncitons.InternalServerError(w, r)
			return
		}

		var requestex string
		requestex =""
		// check if any of the fields are empty
		if strings.TrimSpace(email) == "" || strings.TrimSpace(username) == "" || strings.TrimSpace(password) == "" || strings.Contains(username, "\n") || strings.Contains(email, "\n") || strings.Contains(password, "\n") || strings.Contains(username, " ") || strings.Contains(email, " ") || strings.Contains(password, " ") {
			requestex = "empty"
			tmpl.ExecuteTemplate(w, "register.html", requestex) // requestex sent to html to print message if email/username exists
			return
		}

		if !IsAscii(email) || !IsAscii(password) {
			requestex = "empty"
			tmpl.ExecuteTemplate(w, "register.html", requestex) // requestex sent to html to print message if email/username exists
			return
		}

		// check if username already exists in database
		stmt := "SELECT username FROM User WHERE username = ?"
		row := db.QueryRow(stmt, username) // execute the statment and replaces '?' with 'username' value from form
		var uID string
		err = row.Scan(&uID)
		if err != sql.ErrNoRows || strings.EqualFold(username,"guest"){
			requestex = "user"
		}

		// check if email already exists in database
		stmt = "SELECT email FROM User WHERE email = ?"
		row = db.QueryRow(stmt, email)
		// scan row for username and if found add it to requestex
		err = row.Scan(&uID)
		if err != sql.ErrNoRows || email == "" {
			requestex = requestex + "email"
		}

		// re-renders form with error if username or email already exists
		if requestex != "" {
			tmpl.ExecuteTemplate(w, "register.html", requestex) // requestex sent to html to print message if email/username exists
			return
		}

		if len(strings.TrimSpace(password)) < 8 {
			requestex = "password"
		}

		if requestex == "password" {
			tmpl.ExecuteTemplate(w, "register.html", requestex) // requestex sent to html to print message if email/username exists
			return
		}
		// hash the password
		var hash []byte
		hash, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			fmt.Println("unable to hash password:", err)
			handlerfuncitons.InternalServerError(w, r)
			return
		}
		password = string(hash) // convert hash to string to save it in database

		// Insert new user data into the database
		_, err = db.Exec("INSERT INTO User (email, username, password, gender,pfplink) VALUES (?, ?, ?, ?, ?)", email, username, password, gender,pfp)
		if err != nil {
			fmt.Println("Unable to insert data to the database:", err)
			handlerfuncitons.InternalServerError(w, r)
			return
		}
		sessionToken, err := GenerateSessionToken()
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
			HttpOnly: true,
			Secure:   true,
		})
		http.Redirect(w, r, "/welcome", http.StatusFound)

	} else {
		// Render the form
		tmpl, err := template.ParseFiles("Templates/register.html")
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

func IsAscii(s string) bool {
	for _, c := range s {
		if c > 127 {
			return false
		}
	}
	return true
}

package auth

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	handlerfuncitons "forum/handlers"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type GoogleUserInfo struct {
	Email          string `json:"email"`
	FamilyName     string `json:"family_name"`
	GivenName      string `json:"given_name"`
	ID             string `json:"id"`
	Name           string `json:"name"`
	ProfilePicture string `json:"picture"`
	VerifiedEmail  bool   `json:"verified_email"`
}

func GoogleLoginHandler(w http.ResponseWriter, r *http.Request) {

	state, err := GenerateSessionToken()
	if err != nil {
		fmt.Println("Failed to generate state token")
		handlerfuncitons.InternalServerError(w, r)
		return
	}

	// Store the state in a session or temporary storage (like a cookie)
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   state, // we use this to validate if we made the request
		Expires: time.Now().Add(10 * time.Minute),
		Path:    "/",
	})

	// From google cloud we get the client ID
	clientID := "592535746870-99ugburqi350g8te68lidmhvmp9a3ioh.apps.googleusercontent.com"
	redirectURI := "http://localhost:8080/auth/google/callback"
	authURL := "https://accounts.google.com/o/oauth2/auth"

	u, err := url.Parse(authURL)
	if err != nil {
		http.Error(w, "Failed to parse auth URL", http.StatusInternalServerError)
		return
	}
	q := u.Query()
	q.Set("client_id", clientID)
	q.Set("redirect_uri", redirectURI)
	q.Set("response_type", "code")
	q.Set("scope", "email profile")
	q.Set("state", state)
	u.RawQuery = q.Encode()

	// this actually redirect the user to google authentication website
	http.Redirect(w, r, u.String(), http.StatusTemporaryRedirect)
}

// this is where google responce come to
func GoogleCallbackHandler(w http.ResponseWriter, r *http.Request) {

	// Retrieve the state from the cookie to validate it
	cookie, err := r.Cookie("session_token")
	if err != nil || r.URL.Query().Get("state") != cookie.Value {
		http.Error(w, "State is invalid", http.StatusBadRequest)
		return
	}

	// we get the code google sent from the url
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "No code in the response", http.StatusBadRequest)
		return
	}

	// we exchange the code for the token
	token, err := exchangeCodeForToken(code)
	if err != nil {
		fmt.Println("Error: ", err)
		http.Error(w, "Failed to exchange code for token", http.StatusInternalServerError)
		return
	}

	userInfo, err := fetchGoogleUserInfo(token)
	if err != nil {
		http.Error(w, "Failed to fetch user info", http.StatusInternalServerError)
		return
	}

	// Check if user is already registered database
	isRegistered, err := isEmailRegisteredInDatabase(userInfo.Email)
	if err != nil {
		fmt.Println("Error in google line 100: ", err)
		handlerfuncitons.InternalServerError(w, r)
		return
	}

	// open database because we will use it below
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if isRegistered {
		// login the user normally + redirect to home page
		sessionToken, err := GenerateSessionToken()
		if err != nil {
			handlerfuncitons.InternalServerError(w, r)
			return
		}

		// get the username from the email
		var username string
		err = db.QueryRow("SELECT username FROM User WHERE email = ?", userInfo.Email).Scan(&username)
		if err != nil && err == sql.ErrNoRows {
			fmt.Println("Error: Email not found in the DataBase: ", err)
			handlerfuncitons.InternalServerError(w, r)
			return
		}

		if err := SetSessionToken(db, username, sessionToken); err != nil {
			handlerfuncitons.InternalServerError(w, r)
			return
		}
		// Assuming tmpl is properly initialized
		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   sessionToken,
			Path:    "/",
			Expires: time.Now().Add(24 * time.Hour),
		})
		// if successful then redirect user to main page
		http.Redirect(w, r, "/welcome", http.StatusFound)

	} else {
		// register the user
		// check if username is unique to add the user the the database
		username := userInfo.GivenName + " " + userInfo.FamilyName
		username, err = CheckUniqueUsername(username)
		if err != nil {
			fmt.Println("Error: ", err)
			handlerfuncitons.InternalServerError(w, r)
			return
		}
		pass := make([]byte, 16)
		if _, err := rand.Read(pass); err != nil {
			handlerfuncitons.InternalServerError(w, r)
			return
		}
		password := "GoogleAuth" + string(pass)
		gender := "Google"

		var hash []byte
		hash, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		if err != nil {
			fmt.Println("unable to hash password:", err)
			handlerfuncitons.InternalServerError(w, r)
			return
		}
		password = string(hash) // convert hash to string to save it in database

		// Insert new user data into the database
		_, err = db.Exec("INSERT INTO User (email, username, password, gender,pfplink) VALUES (?, ?, ?, ?,?)", userInfo.Email, username, password, gender,userInfo.ProfilePicture)
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
		// Assuming tmpl is properly initialized
		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   sessionToken,
			Path:    "/",
			Expires: time.Now().Add(24 * time.Hour),
		})
		// if successful then redirect user to main page
		http.Redirect(w, r, "/welcome", http.StatusFound)
	}

}

func exchangeCodeForToken(code string) (string, error) {
	tokenURL := "https://oauth2.googleapis.com/token"
	clientID := "592535746870-99ugburqi350g8te68lidmhvmp9a3ioh.apps.googleusercontent.com"
	clientSecret := "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"
	redirectURI := "http://localhost:8080/auth/google/callback"

	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret) // this is from google cloud website to validate that this is the actual website
	data.Set("code", code)
	data.Set("grant_type", "authorization_code") // what we are going to give in exchange of the info
	data.Set("redirect_uri", redirectURI)

	// pass the data to the http request
	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	// specify the content type in the header
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req) // send the http request
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// save the response
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	// extract the access_token from the resul
	tokenValue := result["access_token"]

	// make sure the token is not nil
	if tokenValue != nil {
		token := tokenValue.(string)
		return token, nil
	}

	// if there is an issue return error
	return "", fmt.Errorf("no access token in response")
}

func fetchGoogleUserInfo(token string) (GoogleUserInfo, error) {
	// url to get the user info
	userInfoURL := "https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token

	var userInfo GoogleUserInfo

	// send a get request to the url to get the user info
	resp, err := http.Get(userInfoURL)
	if err != nil {
		return userInfo, err
	}
	defer resp.Body.Close()

	// save the user info
	json.NewDecoder(resp.Body).Decode(&userInfo)
	return userInfo, nil
}

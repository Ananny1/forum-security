package auth

// import (
// 	//"database/sql"
// 	"crypto/rand"
// 	"database/sql"
// 	"encoding/json"
// 	"fmt"
// 	handlerfuncitons "forum/handlers"
// 	"log"

// 	//"log"
// 	"net/http"
// 	"net/url"
// 	"strings"
// 	"time"

// 	"golang.org/x/crypto/bcrypt"
// )

// type GitHubUserInfo struct {
// 	Login     string `json:"login"`
// 	ID        int    `json:"id"`
// 	Name      string `json:"name"`
// 	AvatarURL string `json:"avatar_url"`
// 	Email     string `json:"email"`
// }

// // GitHubHandler initiates the GitHub OAuth flow
// func GitHubHandler(w http.ResponseWriter, r *http.Request) {
// 	state, err := GenerateSessionToken()
// 	if err != nil {
// 		fmt.Println("err", err)
// 		handlerfuncitons.InternalServerError(w, r)
// 		return
// 	}

// 	http.SetCookie(w, &http.Cookie{
// 		Name:    "session_token",
// 		Value:   state,
// 		Expires: time.Now().Add(24 * time.Hour),
// 		Path:    "/",
// 	})

// 	clientId := "Ov23liQfTrnImNAIfUSC"
// 	redirectURI := "http://localhost:8080/auth/github/callback"
// 	authURL := "https://github.com/login/oauth/authorize"

// 	u, err := url.Parse(authURL)
// 	if err != nil {
// 		http.Error(w, "Failed to parse auth URL", http.StatusInternalServerError)
// 		return
// 	}
// 	q := u.Query()
// 	q.Set("client_id", clientId)
// 	q.Set("redirect_uri", redirectURI)
// 	q.Set("scope", "user:email")
// 	q.Set("state", state)
// 	u.RawQuery = q.Encode()
// 	// http.Redirect(w, r, u.String(), http.StatusFound)

// 	http.Redirect(w, r, u.String(), http.StatusTemporaryRedirect)

// }

// // GitHubCallbackHandler handles the OAuth callback
// func GitHubCallbackHandler(w http.ResponseWriter, r *http.Request) {

// 	//1. get the state and code from the request query
// 	state := r.URL.Query().Get("state")
// 	code := r.URL.Query().Get("code")
// 	if state == "" {
// 		http.Error(w, "Missing state or code", http.StatusBadRequest)
// 		return
// 	}

// 	if code == "" {
// 		http.Error(w, "Missing code", http.StatusBadRequest)
// 		return
// 	}

// 	//1.1 Validate state (not that importnent just extra security)
// 	cookie, err := r.Cookie("session_token")
// 	if err != nil || cookie.Value != state {
// 		http.Error(w, "Invalid state", http.StatusBadRequest)
// 		return
// 	}

// 	// 2.exchange the code for token
// 	token, err := exchangeGitHubCodeForToken(code)
// 	if err != nil {
// 		fmt.Println("Error: ", err)
// 		http.Error(w, "Failed to exchange code for token", http.StatusInternalServerError)
// 		return
// 	}

// 	//3.featch the user infp
// 	userInfo, err := fetchGitHubUserInfo(token)
// 	if err != nil {
// 		http.Error(w, "Failed to fetch user info", http.StatusInternalServerError)
// 		//fmt.Println("Error: ", err)
// 		return
// 	}

// 	GithubisRegistered, err := isEmailRegisteredInDatabase(userInfo.Email)
// 	if err != nil {
// 		fmt.Println("Error in github: ", err)
// 		handlerfuncitons.InternalServerError(w, r)
// 		return
// 	}
// 	db, err := sql.Open("sqlite3", "./forum.db")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer db.Close()

// 	if GithubisRegistered {
// 		sessionToken, err := GenerateSessionToken()
// 		if err != nil {
// 			handlerfuncitons.InternalServerError(w, r)
// 			return
// 		}

// 		// get the username from the email
// 		var username string
// 		err = db.QueryRow("SELECT username FROM User WHERE email = ?", userInfo.Email).Scan(&username)
// 		if err != nil && err == sql.ErrNoRows {
// 			fmt.Println("Error: Email not found in the DataBase: ", err)
// 			handlerfuncitons.InternalServerError(w, r)
// 			return
// 		}

// 		if err := SetSessionToken(db, username, sessionToken); err != nil {
// 			handlerfuncitons.InternalServerError(w, r)
// 			return
// 		}
// 		// Assuming tmpl is properly initialized
// 		http.SetCookie(w, &http.Cookie{
// 			Name:    "session_token",
// 			Value:   sessionToken,
// 			Path:    "/",
// 			Expires: time.Now().Add(24 * time.Hour),
// 		})
// 		// if successful then redirect user to main page
// 		http.Redirect(w, r, "/welcome", http.StatusFound)
// 	} else {
// 		// register the user
// 		// check if username is unique to add the user the the database
// 		username := userInfo.Name
// 		username, err = CheckUniqueUsername(username)
// 		if err != nil {
// 			fmt.Println("Error: ", err)
// 			handlerfuncitons.InternalServerError(w, r)
// 			return
// 		}
// 		pass := make([]byte, 16)
// 		if _, err := rand.Read(pass); err != nil {
// 			handlerfuncitons.InternalServerError(w, r)
// 			return
// 		}
// 		password := "GoogleAuth" + string(pass)
// 		gender := "Github"

// 		var hash []byte
// 		hash, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

// 		if err != nil {
// 			fmt.Println("unable to hash password:", err)
// 			handlerfuncitons.InternalServerError(w, r)
// 			return
// 		}
// 		password = string(hash) // convert hash to string to save it in database		// Insert new user data into the database

// 		_, err = db.Exec("INSERT INTO User (email, username, password, gender,pfplink) VALUES (?, ?, ?, ?,?)", userInfo.Email, username, password, gender,userInfo.AvatarURL)
// 		if err != nil {
// 			fmt.Println("Unable to insert data to the database:", err)
// 			handlerfuncitons.InternalServerError(w, r)
// 			return
// 		}
// 		sessionToken, err := GenerateSessionToken()
// 		if err != nil {
// 			handlerfuncitons.InternalServerError(w, r)
// 			return
// 		}
// 		if err := SetSessionToken(db, username, sessionToken); err != nil {
// 			handlerfuncitons.InternalServerError(w, r)
// 			return
// 		}
// 		// Assuming tmpl is properly initialized
// 		http.SetCookie(w, &http.Cookie{
// 			Name:    "session_token",
// 			Value:   sessionToken,
// 			Path:    "/",
// 			Expires: time.Now().Add(24 * time.Hour),
// 		})
// 		// if successful then redirect user to main page
// 		http.Redirect(w, r, "/welcome", http.StatusFound)
// 	}

// }

// // exchangeGitHubCodeForToken exchanges the authorization code for an access token
// func exchangeGitHubCodeForToken(code string) (string, error) {
// 	tokenURL := "https://github.com/login/oauth/access_token"
// 	clientID := "Ov23liQfTrnImNAIfUSC"
// 	clientSecret := "97acf4a73a6519a6684098381964b73ba4cf6924"
// 	redirectURI := "http://localhost:8080/auth/github/callback"

// 	data := url.Values{}
// 	data.Set("client_id", clientID)
// 	data.Set("client_secret", clientSecret)
// 	data.Set("code", code)
// 	data.Set("redirect_uri", redirectURI)

// 	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
// 	if err != nil {
// 		return "", err
// 	}
// 	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
// 	req.Header.Set("Accept", "application/json")

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return "", err
// 	}
// 	defer resp.Body.Close()

// 	var result map[string]interface{}
// 	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
// 		return "", err
// 	}

// 	tokenValue, alr := result["access_token"].(string)
// 	if !alr {
// 		return "", fmt.Errorf("no access token in response")
// 	}
// 	return tokenValue, nil
// }

// // fetchGitHubUserInfo retrieves the authenticated user's GitHub profile
// func fetchGitHubUserInfo(token string) (GitHubUserInfo, error) {
// 	userInfoURL := "https://api.github.com/user"
// 	emailInfoURL := "https://api.github.com/user/emails"
// 	var userInfo GitHubUserInfo

// 	// Request user info
// 	req, err := http.NewRequest("GET", userInfoURL, nil)
// 	if err != nil {
// 		return userInfo, err
// 	}
// 	req.Header.Set("Authorization", "Bearer "+token)

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return userInfo, err
// 	}
// 	defer resp.Body.Close()

// 	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
// 		return userInfo, err
// 	}

// 	// If no email is present, fetch emails using the /user/emails endpoint
// 	if userInfo.Email == "" {
// 		req, err := http.NewRequest("GET", emailInfoURL, nil)
// 		if err != nil {
// 			return userInfo, err
// 		}
// 		req.Header.Set("Authorization", "Bearer "+token)

// 		resp, err := client.Do(req)
// 		if err != nil {
// 			return userInfo, err
// 		}
// 		defer resp.Body.Close()

// 		var emails []struct {
// 			Email    string `json:"email"`
// 			Primary  bool   `json:"primary"`
// 			Verified bool   `json:"verified"`
// 		}

// 		if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
// 			return userInfo, err
// 		}

// 		// Set the user's email to their primary email if available
// 		for _, email := range emails {
// 			if email.Primary && email.Verified {
// 				userInfo.Email = email.Email
// 				break
// 			}
// 		}
// 	}

// 	return userInfo, nil
// }

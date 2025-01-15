package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/rvaidun/svmail/mydatabase"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	googleouath2 "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

// Scopes: OAuth 2.0 scopes provide a way to limit the amount of access that is granted to an access token.
var googleOauthConfig = &oauth2.Config{
	RedirectURL:  "http://localhost:8000/auth/google/callback",
	ClientID:     "347014314619-88tqtgk7a71g2c3ra5bg3ombejr2sgi6.apps.googleusercontent.com",
	ClientSecret: "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",
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://mail.google.com/"},
	Endpoint:     google.Endpoint,
}

type Session struct {
	User      mydatabase.User
	CSRFToken string
}

var SESSIONS = make(map[string]Session)

func oauthGoogleLogin(w http.ResponseWriter, r *http.Request) {

	// Create oauthState cookie
	oauthState := generateStateOauthCookie(w)

	/*
		AuthCodeURL receive state that is a token to protect the user from CSRF attacks. You must always provide a non-empty string and
		validate that it matches the the state query parameter on your redirect callback.
	*/
	u := googleOauthConfig.AuthCodeURL(oauthState)
	http.Redirect(w, r, u, http.StatusTemporaryRedirect)
}

func oauthGoogleCallback(w http.ResponseWriter, r *http.Request) {
	// Read oauthState from Cookie
	oauthState, _ := r.Cookie("oauthstate")

	if r.FormValue("state") != oauthState.Value {
		log.Println("invalid oauth google state")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	token, err := getTokenFromCode(r.FormValue("code"))
	if err != nil {
		log.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	dbConn, err := mydatabase.CreateConn()
	if err != nil {
		log.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	data, err := GetUserDataFromGoogle(token)
	if err != nil {
		log.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	userExists, _ := mydatabase.UserExists(dbConn, data.Email)
	if !userExists {
		err = mydatabase.CreateUser(dbConn, token, data.Email)
		if err != nil {
			log.Println(err.Error())
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

	}
	fmt.Printf("User: %v in database\n", data)
	// Create session ID and csrf token
	sessionID := generateSessionID()
	csrfToken := generateSessionID()
	SESSIONS[sessionID] = Session{User: mydatabase.User{Token: *token, Email: data.Email}, CSRFToken: csrfToken}

	// set session id in cookies
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Expires:  time.Now().Add(60 * time.Hour),
		Secure:   true,
		HttpOnly: true,
		Path:     "/",
	})

	// return csrf token in the body
	var jsonResponse = map[string]string{"csrf_token": csrfToken}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(jsonResponse)

}

func oauthGoogleLogout(w http.ResponseWriter, r *http.Request) {
	sessionID, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	delete(SESSIONS, sessionID.Value)
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Expires:  time.Now(),
		Secure:   true,
		HttpOnly: true,
	})
	w.Write([]byte("Logged out"))
}

func generateSessionID() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func generateStateOauthCookie(w http.ResponseWriter) string {
	var expiration = time.Now().Add(20 * time.Minute)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(w, &cookie)

	return state
}

func getTokenFromCode(code string) (*oauth2.Token, error) {
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange wrong: %s", err.Error())
	}
	return token, nil
}

func GetUserDataFromGoogle(token *oauth2.Token) (*googleouath2.Userinfo, error) {
	httpClient := googleOauthConfig.Client(context.Background(), token)
	// service, err := googleouath2.New(httpClient)
	service, err := googleouath2.NewService(context.Background(), option.WithHTTPClient(httpClient))
	if err != nil {
		return nil, fmt.Errorf("unable to create google service: %v", err)
	}
	userinfo, err := service.Userinfo.Get().Do()
	if err != nil {
		return nil, fmt.Errorf("unable to get userinfo: %v", err)
	}
	return userinfo, nil
}

func GetGoogleOauthClient() *googleouath2.Service {
	httpClient := googleOauthConfig.Client(context.Background(), nil)
	service, err := googleouath2.NewService(context.Background(), option.WithHTTPClient(httpClient))
	if err != nil {
		log.Fatalf("unable to create google service: %v", err)
	}
	return service
}

type contextKey string

const UserKey = contextKey("user")

func AuthenticatedMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "moz-extension://0d05fa30-b941-4dad-9abd-9fadee86fbe8")
			// remove Content-Type from preflight request
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Cookie")

			return
		}
		sessionID, err := r.Cookie("session_id")
		if err != nil {
			fmt.Println("No session id cookie")
			w.Header().Set("Access-Control-Allow-Origin", "moz-extension://0d05fa30-b941-4dad-9abd-9fadee86fbe8")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		session, ok := SESSIONS[sessionID.Value]
		if !ok {
			fmt.Println("No session found")
			w.Header().Set("Access-Control-Allow-Origin", "moz-extension://0d05fa30-b941-4dad-9abd-9fadee86fbe8")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// check if CSRF token is valid
		// csrfToken := r.Header.Get("X-CSRF-Token")
		// if csrfToken != session.CSRFToken {
		// 	fmt.Println("Invalid CSRF token")
		// 	http.Error(w, "Unauthorized", http.StatusUnauthorized)
		// 	return
		// }
		fmt.Println("Authenticated the following email address: ", session.User.Email)
		w.Header().Set("Access-Control-Allow-Origin", "moz-extension://0d05fa30-b941-4dad-9abd-9fadee86fbe8")
		ctx := context.WithValue(r.Context(), UserKey, &session.User)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

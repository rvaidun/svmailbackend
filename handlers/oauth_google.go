package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
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
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
	Endpoint:     google.Endpoint,
}

const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

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
	data, err := getUserDataFromGoogle(token)
	if err != nil {
		log.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	err = mydatabase.CreateUserWithToken(dbConn, token, data.Email)
	if err != nil {
		log.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	fmt.Printf("User: %v in database\n", data)

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

func getUserDataFromGoogle(token *oauth2.Token) (*googleouath2.Userinfo, error) {
	// Use code to get token and get user info from Google.
	// response, err := http.Get(oauthGoogleUrlAPI + token.AccessToken)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	// }
	// defer response.Body.Close()
	// contents, err := io.ReadAll(response.Body)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed read response: %s", err.Error())
	// }
	// // parse []byte to json

	// userInfo := UserInfo{}
	// err = json.Unmarshal(contents, &userInfo)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed unmarshal response: %s", err.Error())
	// }

	// return &userInfo, nil
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

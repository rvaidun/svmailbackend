package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rvaidun/svmail/mydatabase"
)

func userInfo(w http.ResponseWriter, r *http.Request) {
	// We get the user from the context
	w.Header().Set("Access-Control-Allow-Origin", "*")
	user := r.Context().Value(UserKey).(*mydatabase.User)
	// We write the user information to the response
	// get the user data from google
	googleUserData, err := GetUserDataFromGoogle(&user.Token)
	if err != nil {
		fmt.Printf("Error getting user data from google: %v\n", err)
		return
	}
	// write the user data to the response as json

	jsonData, err := json.Marshal(googleUserData)
	if err != nil {
		fmt.Printf("Error marshalling user data: %v\n", err)
		http.Error(w, "Error marshalling user data", http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)

}

func getUserEmail(w http.ResponseWriter, r *http.Request) {
	// We get the user from the context
	user := r.Context().Value(UserKey).(*mydatabase.User)
	// We write the user information to the response
	emailId := r.URL.Query().Get("email")
	if emailId == "" {
		w.Write([]byte(fmt.Sprintf("Hello, %s\n", user.Email)))
		return
	}
	emailData, err := GetEmailDataFromGoogle(&user.Token, emailId)
	if err != nil {
		fmt.Printf("Error getting email data: %v\n", err)
		w.Write([]byte(fmt.Sprintf("Error getting email data: %v\n", err)))
		return
	}
	// write the email data to the response
	w.Write([]byte(fmt.Sprintf("Email data: %v\n", emailData)))

}

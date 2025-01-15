package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rvaidun/svmail/mydatabase"
)

func handleGetTrackedEmails(w http.ResponseWriter, r *http.Request) {
	// We get the user from the context
	user := r.Context().Value(UserKey).(*mydatabase.User)
	// We write the user information to the response
	w.Write([]byte(fmt.Sprintf("Hello, %s\n", user.Email)))
	// get the tracked emails from the database
	conn, err := mydatabase.CreateConn()
	if err != nil {
		fmt.Printf("Error creating connection: %v\n", err)
		return
	}
	// get thread ids from the body
	var thread_ids []string
	err = json.NewDecoder(r.Body).Decode(&thread_ids)
	if err != nil {
		fmt.Printf("Error decoding thread ids: %v\n", err)
		return
	}
	// get the tracked emails from the database
	trackedEmails, err := mydatabase.GetTrackedEmails(conn, thread_ids, user.Email)
	if err != nil {
		fmt.Printf("Error getting tracked emails: %v\n", err)
		return
	}
	// write the tracked emails to the response as json
	jsonResponse, err := json.Marshal(trackedEmails)
	if err != nil {
		fmt.Printf("Error marshalling tracked emails: %v\n", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func handlePostTrackedEmail(w http.ResponseWriter, r *http.Request) {
	// We get the user from the context
	user := r.Context().Value(UserKey).(*mydatabase.User)
	// We write the user information to the response
	w.Write([]byte(fmt.Sprintf("Hello, %s\n", user.Email)))
	// get the tracked email from the body
	var messageID string
	var threadID string
	err := json.NewDecoder(r.Body).Decode(&struct {
		MessageID *string `json:"message_id"`
		ThreadID  *string `json:"thread_id"`
	}{&messageID, &threadID})
	if err != nil {
		fmt.Printf("Error decoding tracked email: %v\n", err)
		return
	}
	// create the tracked email in the database
	conn, err := mydatabase.CreateConn()
	if err != nil {
		fmt.Printf("Error creating connection: %v\n", err)
		return
	}
	err = mydatabase.CreateTrackedEmail(conn, messageID, threadID, user.Email)
	if err != nil {
		fmt.Printf("Error creating tracked email: %v\n", err)
		return
	}
	w.Write([]byte("Tracked email created"))
}

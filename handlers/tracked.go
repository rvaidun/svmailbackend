package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

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
	// get thread ids from the query string in the request URI
	query := r.URL.Query()
	threadIDs := query["thread_id"]

	fmt.Printf("Thread ids: %v\n", threadIDs)

	// get the tracked emails from the database
	trackedEmails, err := mydatabase.GetTrackedEmails(conn, threadIDs, user.Email)
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

func handleViewCount(w http.ResponseWriter, r *http.Request) {
	// get the id from the request
	messageID := r.PathValue("message_id")

	curUnixTime := time.Now().Unix()

	ip := r.RemoteAddr

	conn, err := mydatabase.CreateConn()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	err = mydatabase.CreateView(conn, messageID, curUnixTime, ip)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// return a 404 page
	http.Error(w, "Not found!", http.StatusNotFound)

}

func handleGetViewsForMessage(w http.ResponseWriter, r *http.Request) {
	// get the id from the request
	messageID := r.PathValue("message_id")

	conn, err := mydatabase.CreateConn()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	views, err := mydatabase.GetViews(conn, messageID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	jsonResponse, err := json.Marshal(views)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

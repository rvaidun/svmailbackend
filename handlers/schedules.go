package handlers

import (
	"fmt"
	"net/http"
)

func getScheduledEmails(w http.ResponseWriter, r *http.Request) {
	// We get the user from the context
	user := r.Context().Value(UserKey).(*mydatabase.User)
	// We write the user information to the response
	w.Write([]byte(fmt.Sprintf("Hello, %s\n", user.Email)))
	// get the scheduled emails from the database
	scheduledEmails, err := mydatabase.GetScheduledEmails(user.Email)
	if err != nil {
		fmt.Printf("Error getting scheduled emails: %v\n", err)
		return
	}
	// write the scheduled emails to the response
	w.Write([]byte(fmt.Sprintf("Scheduled emails: %v\n", scheduledEmails)))
}

package handlers

import (
	"net/http"
)

func methodHelper(w http.ResponseWriter, r *http.Request, getHandler http.HandlerFunc, postHandler http.HandlerFunc) {
	if (r.Method == "GET" || r.Method == "") && getHandler != nil {
		getHandler(w, r)
	} else if r.Method == "POST" && postHandler != nil {
		postHandler(w, r)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
func New() http.Handler {
	mux := http.NewServeMux()
	// Root
	mux.Handle("/", http.FileServer(http.Dir("templates/")))

	// OauthGoogle
	mux.HandleFunc("/auth/google/login", oauthGoogleLogin)
	mux.HandleFunc("/auth/google/callback", oauthGoogleCallback)

	// different routes for GET and POST of /tracked
	mux.HandleFunc("/tracked", func(w http.ResponseWriter, r *http.Request) {
		methodHelper(w, r, handleGetTrackedEmails, handlePostTrackedEmail)
	})

	return mux
}

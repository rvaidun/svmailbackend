package handlers

import (
	"net/http"
)

func methodHelper(w http.ResponseWriter, r *http.Request, getHandler http.Handler, postHandler http.Handler) {
	if (r.Method == "GET" || r.Method == "") && getHandler != nil {
		getHandler.ServeHTTP(w, r)
	} else if r.Method == "POST" && postHandler != nil {
		postHandler.ServeHTTP(w, r)
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
		methodHelper(w, r, AuthenticatedMiddleware(http.HandlerFunc(handleGetTrackedEmails)), AuthenticatedMiddleware(http.HandlerFunc(handlePostTrackedEmail)))
	})

	mux.Handle("/test", AuthenticatedMiddleware(http.HandlerFunc(getUserEmail)))

	// handler for viewing scheduled emails
	mux.HandleFunc("/imgs/{message_id}.jpg", handleViewCount)

	return mux
}

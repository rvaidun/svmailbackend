package handlers

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/rvaidun/svmail/mydatabase"
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

type IndexPageData struct {
	IsAuthenticated bool
	User            mydatabase.User
}

func New() http.Handler {
	mux := http.NewServeMux()
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	// Root and pass template to it
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		sessionID, err := r.Cookie("session_id")
		if err != nil {
			tmpl.Execute(w, IndexPageData{IsAuthenticated: false})
			return
		}
		session, ok := SESSIONS[sessionID.Value]
		if !ok {
			fmt.Println("No session found")
			tmpl.Execute(w, IndexPageData{IsAuthenticated: false})
			return
		}
		tmpl.Execute(w, IndexPageData{IsAuthenticated: true, User: session.User})

	})

	// OauthGoogle
	mux.HandleFunc("/auth/google/login", oauthGoogleLogin)
	mux.HandleFunc("/auth/google/callback", oauthGoogleCallback)
	mux.HandleFunc("/auth/google/logout", oauthGoogleLogout)

	// different routes for GET and POST of /tracked
	mux.HandleFunc("/tracked", func(w http.ResponseWriter, r *http.Request) {
		methodHelper(w, r, AuthenticatedMiddleware(http.HandlerFunc(handleGetTrackedEmails)), AuthenticatedMiddleware(http.HandlerFunc(handlePostTrackedEmail)))
	})

	mux.Handle("/test", AuthenticatedMiddleware(http.HandlerFunc(getUserEmail)))
	mux.Handle("/userinfo", AuthenticatedMiddleware(http.HandlerFunc(userInfo)))

	// handler for viewing scheduled emails
	mux.HandleFunc("/imgs/{message_id}.jpg", handleViewCount)

	return mux
}

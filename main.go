package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/rvaidun/svmail/handlers"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func main() {
	godotenv.Load()

	server := &http.Server{
		Addr:    ":" + os.Getenv("APPLICATION_PORT"),
		Handler: handlers.New(),
	}
	handlers.GoogleOauthConfig = &oauth2.Config{
		RedirectURL:  "https://" + os.Getenv("APPLICATION_HOST") + "/auth/google/callback",
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://mail.google.com/"},
		Endpoint:     google.Endpoint,
	}
	fmt.Printf("DB_PASSWORD: %s\n", os.Getenv("DB_PASSWORD"))

	fmt.Printf("Starting HTTP Server. Listening at %q", server.Addr)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		fmt.Printf("%v", err)
	} else {
		fmt.Println("Server closed!")
	}
}

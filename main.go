package main

import (
	"fmt"
	"net/http"

	"github.com/rvaidun/svmail/handlers"
)

func main() {
	server := &http.Server{
		Addr:    ":8000",
		Handler: handlers.New(),
	}

	fmt.Printf("Starting HTTP Server. Listening at %q", server.Addr)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		fmt.Printf("%v", err)
	} else {
		fmt.Println("Server closed!")
	}
}

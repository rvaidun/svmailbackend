package handlers

import (
	"context"
	"fmt"

	"golang.org/x/oauth2"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

// NewGmailService initializes a new Gmail service client

func GetEmailDataFromGoogle(token *oauth2.Token, emailId string) (*gmail.Message, error) {
	httpClient := googleOauthConfig.Client(context.Background(), token)
	// service, err := googleouath2.New(httpClient)
	service, err := gmail.NewService(context.Background(), option.WithHTTPClient(httpClient))
	if err != nil {
		return nil, fmt.Errorf("unable to create gmail service: %v", err)
	}
	// get the email data
	emailData, err := service.Users.Messages.Get("me", emailId).Do()
	if err != nil {
		return nil, fmt.Errorf("unable to get email data: %v", err)
	}
	return emailData, nil
}

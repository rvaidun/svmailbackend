package mydatabase

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"golang.org/x/oauth2"
)

// creates a user with a oauth2 token in the database
func CreateUserWithToken(conn *pgx.Conn, token *oauth2.Token, email string) error {
	unixTime := token.Expiry.Unix()
	_, err := conn.Exec(context.Background(), "INSERT INTO users (access_token, token_type, refresh_token, expiry, email) VALUES ($1, $2, $3, $4, $5)", token.AccessToken, token.TokenType, token.RefreshToken, unixTime, email)
	if err != nil {
		fmt.Printf("Error creating user with token: %v\n", err)
		return err
	}
	return nil

}

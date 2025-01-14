package mydatabase

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"golang.org/x/oauth2"
)

type User struct {
	Token oauth2.Token
	Email string
}

var DatabaseConn *pgx.Conn = nil

func CreateConn() (*pgx.Conn, error) {
	if DatabaseConn == nil {
		conn, err := pgx.Connect(context.Background(), "postgres://postgres@localhost:5432/svmail")
		if err != nil {
			return nil, fmt.Errorf("unable to connect to database: %v", err)
		}
		DatabaseConn = conn
	}
	return DatabaseConn, nil
}

package mydatabase

import (
	"context"
	"fmt"
	"os"

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
		dbConnString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
		conn, err := pgx.Connect(context.Background(), dbConnString)
		if err != nil {
			return nil, fmt.Errorf("unable to connect to database: %v", err)
		}
		DatabaseConn = conn
	}
	return DatabaseConn, nil
}

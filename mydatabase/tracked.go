package mydatabase

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// returns map of thread id to list of message ids for tracked emails
// create map type
func GetTrackedEmails(conn *pgx.Conn, thread_ids []string, username string) (map[string][]string, error) {
	rows, err := conn.Query(context.Background(), "SELECT thread_id, message_id FROM tracked_emails WHERE thread_id = ANY($1) AND username = $2", thread_ids, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trackedEmails = make(map[string][]string)
	for rows.Next() {
		var threadID string
		var messageID string
		err := rows.Scan(&threadID, &messageID)
		if err != nil {
			return nil, err
		}
		// if the key does not exist, create it
		if _, ok := trackedEmails[threadID]; !ok {
			trackedEmails[threadID] = []string{}
		}
		trackedEmails[threadID] = append(trackedEmails[threadID], messageID)
	}
	return trackedEmails, nil
}

func CreateTrackedEmail(conn *pgx.Conn, threadID string, messageID string, username string) error {
	_, err := conn.Exec(context.Background(), "INSERT INTO tracked_emails (thread_id, message_id, username) VALUES ($1, $2, $3)", threadID, messageID, username)
	if err != nil {
		return err
	}
	return nil
}

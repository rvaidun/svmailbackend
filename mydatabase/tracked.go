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

type TrackedEmail struct {
	ThreadID  string `json:"thread_id"`
	MessageID string `json:"message_id"`
	Username  string `json:"username"`
}

func CreateView(conn *pgx.Conn, messageID string, time int64, ip string) error {
	_, err := conn.Exec(context.Background(), "INSERT INTO email_views (message_id, time, ip) VALUES ($1, $2, $3)", messageID, time, ip)
	if err != nil {
		return err
	}
	return nil
}

type EmailView struct {
	MessageID string `json:"message_id"`
	ViewTime  int64  `json:"viewed_time"`
	IP        string `json:"ip"`
}

func GetViews(conn *pgx.Conn, messageID string) ([]EmailView, error) {
	rows, err := conn.Query(context.Background(), "SELECT message_id, time, ip FROM email_views WHERE message_id = $1", messageID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var views []EmailView
	for rows.Next() {
		var view EmailView
		err := rows.Scan(&view.MessageID, &view.ViewTime, &view.IP)
		if err != nil {
			return nil, err
		}
		views = append(views, view)
	}
	return views, nil
}

func GetLatestView(conn *pgx.Conn, messageIDs []string) ([]EmailView, error) {
	// get the latest view for each message id
	rows, err := conn.Query(context.Background(), "SELECT DISTINCT ON (message_id) message_id, time, ip FROM email_views WHERE message_id = ANY($1) ORDER BY message_id, time DESC", messageIDs)
	if err != nil {
		return nil, err
	}
	// add all messageIDs to a Set
	messageIDSET := make(map[string]bool)
	defer rows.Close()
	var views []EmailView
	for rows.Next() {
		var view EmailView
		err := rows.Scan(&view.MessageID, &view.ViewTime, &view.IP)
		if err != nil {
			return nil, err
		}
		views = append(views, view)
		messageIDSET[view.MessageID] = true
	}
	// for any message id that does not have a view, add an empty view

	for _, messageID := range messageIDs {
		if _, ok := messageIDSET[messageID]; !ok {
			views = append(views, EmailView{MessageID: messageID})
		}
	}
	return views, nil
}

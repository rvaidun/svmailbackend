package mydatabase

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type ScheduledEmail struct {
	Email         string
	ScheduledTime int64
	ReadReceipt   bool
}

func GetScheduledEmails(conn *pgx.Conn, email []string) ([]ScheduledEmail, error) {
	rows, err := conn.Query(context.Background(), "SELECT email, scheduled_time, read_receipt FROM scheduled_emails WHERE email = ANY($1)", email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scheduledEmails []ScheduledEmail
	for rows.Next() {
		var scheduledEmail ScheduledEmail
		err := rows.Scan(&scheduledEmail.Email, &scheduledEmail.ScheduledTime, &scheduledEmail.ReadReceipt)
		if err != nil {
			return nil, err
		}
		scheduledEmails = append(scheduledEmails, scheduledEmail)
	}
	return scheduledEmails, nil
}

func CreateScheduledEmail(conn *pgx.Conn, email string, scheduledTime int64, readReceipt bool) error {
	_, err := conn.Exec(context.Background(), "INSERT INTO scheduled_emails (email, scheduled_time, read_receipt) VALUES ($1, $2, $3)", email, scheduledTime, readReceipt)
	if err != nil {
		return err
	}
	return nil
}

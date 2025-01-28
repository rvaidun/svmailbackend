# svmailbackend
This is a simple backend for a web extension that allows users to see when their emails are opened by the recepient. The backend also uses gmail API to schedule sending emails at a later time.

The backend is written in Go and uses the default net/http package to handle requests. The backend uses a Postgres database to store user data and email data.

## Debugging
```bash
go run main.go
```
## Building
- Build the binary `go build -o svmail`
- Build the binary for linux `GOOS=linux GOARCH=amd64 go build -o svmail-linux`
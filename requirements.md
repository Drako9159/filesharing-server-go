# Initalize module
go mod init server-go

# For run and probe
go run server.go

## For compile for windows stay on windows
go build -o server.exe
## For compile for windows stay on another system
GOOS=windows GOARCH=amd64 go build -o server.exe

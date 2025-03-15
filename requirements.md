# Initalize module
go mod init server-go

# For run and probe
go run server.go

## For compile for windows stay on windows
go build -o server.exe
## For compile for windows stay on another system
GOOS=windows GOARCH=amd64 go build -o server.exe

# For compile info
go install github.com/akavel/rsrc@latest

# Generate .syso
$(go env GOPATH)/bin/rsrc -ico home_server.ico -manifest app.manifest -o rsrc.syso
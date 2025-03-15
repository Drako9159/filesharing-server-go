# Simple File Sharing Server

This project is a lightweight file-sharing server designed to share files from a local folder over a network. It also provides the functionality to upload files, making it a quick and easy solution for transferring files between devices.

## Features

- Share files from a local folder over the network.
- Upload files to the server for quick file transfers.

## Getting Started

### Prerequisites

Make sure you have [Go](https://go.dev/doc/install) installed. You can install it using:

```bash
# Download golang
wget https://go.dev/dl/go1.24.1.linux-amd64.tar.gz

# Remove previous and unpack
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.24.1.linux-amd64.tar.gz

# Add PATH
export PATH=$PATH:/usr/local/go/bin
```
### Running the Server

1. CLone the repository

```bash
git clone https://github.com/Drako9159/filesharing-server-go.git
```
2. Run the server

```bash
go run server.go
```
4. Open your web browser and go to http://localhost:8080/ to access the file-sharing server

## Packaging the Application
If you want to package your application into a single executable, you can use go builder.

1. Install Go

```bash
# Generate .syso
$(go env GOPATH)/bin/rsrc -ico home_server.ico -manifest app.manifest -o rsrc.syso
```

2. Package your application

```bash
GOOS=windows GOARCH=amd64 go build -o server.exe
```

For a complete list of dependencies, refer to the requirements.md file.

#### Additional Notes
- Make sure to customize the application to suit your specific needs.
Security considerations: This application is intended for local network use. If used in a public environment, consider implementing security measures such as authentication and encryption.
- Feel free to explore the code, modify it, and adapt it to your requirements! If you have any questions or suggestions, feel free to open an issue or contribute to the project.
# 📂 My Local Cloud Go - v2.0

A high-performance, lightweight, and secure file sharing server written in Go. This version has been completely redesigned to provide a professional user experience with minimal resource consumption.

## 🚀 What's New in v2.0

- 🌑 **Dark Mode UI**: Sleek, modern, and eye-friendly interface.
- 🔐 **Enhanced Security**: Added password protection with session-based cookies.
- ⚡ **Zero-RAM Streaming**: Uploads are now streamed directly to disk, allowing you to handle files of 10GB+ even on low-spec hardware.
- 🔍 **Instant Search**: Find files in real-time with the new integrated search bar.
- 📱 **Mobile-First Design**: Completely responsive grid-based layout that won't break on small screens.
- 🗑️ **Safe Deletion**: Integrated file removal with a confirmation safety lock.
- 📶 **Network Optimization**: Removed write timeouts to allow stable downloads of massive files over WiFi.

## ✨ Features

- 📂 **Auto-Listing**: Automatically serves files from the execution directory.
- ⬇️ **Smart Downloads**: Supports browser download managers and resuming.
- 📤 **Ultra-Fast Uploads**: Real-time progress bar for tracking status.
- 🛡️ **System Protection**: Automatically hides and protects sensitive files (`main.go`, `go.mod`, `.exe`, etc.).
- 📦 **Self-Contained**: Static files and templates are embedded into the binary.

## 🛠️ Build & Installation

### Prerequisites
- Go 1.20 or higher.

### Run in Development
```bash
go run main.go
```

### Compile for Windows (.exe) from Linux
```bash
GOOS=windows GOARCH=amd64 go build -o cloud_server.exe main.go
```

### Compile for Linux
```bash
go build -o cloud_server main.go
```

## ⚙️ Configuration

Open `main.go` to customize the following:

- **PASSWORD**: Change the access code (default: `1234`).
- **PORT**: Change the server port (default: `8080`).
- **maxUploadSize**: Set your preferred upload limit (default: 10GB).

## 📝 Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/login` | Access authentication page |
| GET | `/` | Main dashboard (requires login) |
| POST | `/upload` | Streaming upload handler |
| GET | `/download/*` | File streaming download |
| POST | `/delete` | Secure file removal |

## ⚠️ Safety Note

This server is designed for **trusted local networks (LAN)**. It is highly efficient for home or office use, but should not be exposed to the public internet without an additional reverse proxy (like Nginx) and HTTPS encryption.

---
Built with ⚡ **Go** for speed and efficiency.

package main

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	//"io/fs"
	//"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/mdp/qrterminal/v3"
)

//go:embed templates/* static/*
var embeddedFiles embed.FS

type FileInfo struct {
	Name string
	Size string
}

const (
	PASSWORD       = "1234" // CAMBIA TU CONTRASEÑA AQUÍ
	SESSION_COOKIE = "auth_token"
	maxUploadSize  = 10 << 30
)

func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit { return fmt.Sprintf("%d B", bytes) }
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func checkAuth(r *http.Request) bool {
	cookie, err := r.Cookie(SESSION_COOKIE)
	return err == nil && cookie.Value == "authorized"
}

func listFiles() ([]FileInfo, error) {
	files, err := os.ReadDir(".")
	if err != nil { return nil, err }
	var list []FileInfo
	for _, f := range files {
		if f.IsDir() { continue }
		name := f.Name()
		if name == "main.go" || name == "go.mod" || strings.HasSuffix(name, ".ico") || strings.HasSuffix(name, ".rc") || strings.HasSuffix(name, ".exe") || strings.HasSuffix(name, ".res") {
			continue
		}
		info, _ := f.Info()
		list = append(list, FileInfo{Name: name, Size: formatSize(info.Size())})
	}
	return list, nil
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl, _ := template.ParseFS(embeddedFiles, "templates/login.html")
		tmpl.Execute(w, nil)
		return
	}
	
	if r.FormValue("password") == PASSWORD {
		http.SetCookie(w, &http.Cookie{
			Name:     SESSION_COOKIE,
			Value:    "authorized",
			Path:     "/",
			HttpOnly: true,
		})
		http.Redirect(w, r, "/", 303)
	} else {
		http.Redirect(w, r, "/login", 303)
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     SESSION_COOKIE,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
	http.Redirect(w, r, "/login", 303)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if !checkAuth(r) { http.Redirect(w, r, "/login", 303); return }
	if r.URL.Path != "/" { http.NotFound(w, r); return }
	files, _ := listFiles()
	tmpl, _ := template.ParseFS(embeddedFiles, "templates/index.html")
	tmpl.Execute(w, files)
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	if !checkAuth(r) { http.Error(w, "Unauthorized", 401); return }
	fileName := strings.TrimPrefix(r.URL.Path, "/download/")
	if fileName == "" || strings.Contains(fileName, "..") {
		http.Error(w, "Invalid file", 400); return
	}
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeFile(w, r, fileName)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if !checkAuth(r) || r.Method != "POST" { return }
	reader, _ := r.MultipartReader()
	for {
		part, err := reader.NextPart()
		if err == io.EOF { break }
		if part.FileName() == "" { continue }
		dst, _ := os.Create(part.FileName())
		io.Copy(dst, part)
		dst.Close()
	}
	http.Redirect(w, r, "/", 303)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	if !checkAuth(r) || r.Method != "POST" { return }
	name := r.FormValue("filename")
	if name == "" || strings.Contains(name, "..") { return }
	os.Remove(name)
	http.Redirect(w, r, "/", 303)
}

func getLocalIP() string {
	addrs, _ := net.InterfaceAddrs()
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ip4 := ipnet.IP.To4(); ip4 != nil {
				if strings.HasPrefix(ip4.String(), "192.168.") { return ip4.String() }
			}
		}
	}
	return "localhost"
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/download/", downloadHandler)
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/delete", deleteHandler)

	port := "8080"
	localIP := getLocalIP()
	serverURL := fmt.Sprintf("http://%s:%s", localIP, port)
	fmt.Println("-------------------------------------------")
	fmt.Println("🚀 GO SECURE SERVER STARTED")
	fmt.Printf("🌐 Network: %s\n", serverURL)
	fmt.Println("📱 Escanea el código QR para acceder:")
	qrterminal.GenerateHalfBlock(serverURL, qrterminal.L, os.Stdout)
	fmt.Println("-------------------------------------------")

	srv := &http.Server{
		Addr:         ":" + port,
		WriteTimeout: 0,
		ReadTimeout:  0,
		IdleTimeout:  60 * time.Second,
	}
	srv.ListenAndServe()
}

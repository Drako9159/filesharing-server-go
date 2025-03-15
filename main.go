package main

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

var tmpl_old = `
<!DOCTYPE html>
<html>
<head>
    <title>File Server</title>
</head>
<body>
    <h1>Files in Directory</h1>
    <ul>
        {{range .}}
        <li><a href="/download/{{.Name}}">{{.Name}} ({{.Size}} KB)</a></li>
        {{end}}
    </ul>
</body>
</html>`

//go:embed templates/* static/*
var embeddedFiles embed.FS

type FileInfo struct {
	Name string
	Size string
}

func listFiles() ([]FileInfo, error) {
	files, err := os.ReadDir(".")
	if err != nil {
		return nil, err
	}

	var fileInfoList []FileInfo
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		info, err := file.Info()
		if err != nil {
			continue
		}
		size := strconv.FormatFloat(float64(info.Size())/1024, 'f', 2, 64)
		fileInfoList = append(fileInfoList, FileInfo{Name: file.Name(), Size: size})
	}
	return fileInfoList, nil
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	files, err := listFiles()
	if err != nil {
		http.Error(w, "Failed to list files", http.StatusInternalServerError)
		return
	}

	// t := template.Must(template.New("index").Parse(tmpl))
	// t.Execute(w, files)
	tmpl, err := template.ParseFS(embeddedFiles, "templates/index.html")
	if err != nil {
		http.Error(w, "Failed to load template", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, files)
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Path[len("/download/"):]
	filePath := filepath.Join(".", fileName)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	http.ServeFile(w, r, filePath)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	dst, err := os.Create(header.Filename)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func staticHandlerOld(w http.ResponseWriter, r *http.Request) {
	file, err := embeddedFiles.ReadFile("static" + r.URL.Path)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	w.Write(file)
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	path := "static" + r.URL.Path[len("/static/"):]
	file, err := embeddedFiles.ReadFile(path)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/css") // Indicar que es un archivo CSS
	w.Write(file)
}

// getLocalIP obtiene la IP local de la máquina en la red
func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "localhost"
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String()
			}
		}
	}
	return "localhost"
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/download/", downloadHandler)
	http.HandleFunc("/upload", uploadHandler)
	// http.HandleFunc("/static/", staticHandlerOld)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	localIP := getLocalIP()
	port := "8080"

	fmt.Printf("Serving on:\n")
	fmt.Printf("➡ Local:   http://localhost:%s\n", port)
	fmt.Printf("➡ Network: http://%s:%s\n", localIP, port)

	http.ListenAndServe(":8080", nil)
}

package main

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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

func staticHandler(w http.ResponseWriter, r *http.Request) {
	// La ruta correcta debe ser así
	path := filepath.Join("static", r.URL.Path[len("/static/"):])

	file, err := embeddedFiles.ReadFile(path)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// Establecer el tipo de contenido basado en la extensión del archivo
	ext := filepath.Ext(path)
	switch ext {
	case ".css":
		w.Header().Set("Content-Type", "text/css")
	case ".js":
		w.Header().Set("Content-Type", "application/javascript")
	case ".png":
		w.Header().Set("Content-Type", "image/png")
	case ".jpg", ".jpeg":
		w.Header().Set("Content-Type", "image/jpeg")
	}

	w.Write(file)
}

func staticHandlerOld(w http.ResponseWriter, r *http.Request) {
	path := "static" + r.URL.Path[len("/static/"):]
	file, err := embeddedFiles.ReadFile(path)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/css") // Indicar que es un archivo CSS
	w.Write(file)
}

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "localhost"
	}

	// Primero intentamos encontrar direcciones IPv4 no locales
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipv4 := ipNet.IP.To4(); ipv4 != nil {
				// Descartamos direcciones de Docker, VirtualBox y otras interfaces virtuales
				if !strings.HasPrefix(ipv4.String(), "172.") &&
					!strings.HasPrefix(ipv4.String(), "10.") &&
					!strings.HasPrefix(ipv4.String(), "192.168.56.") {
					return ipv4.String()
				}
			}
		}
	}

	// Si no encontramos ninguna adecuada, intentamos con cualquier IPv4 privada
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipv4 := ipNet.IP.To4(); ipv4 != nil {
				// Buscamos direcciones de red local típicas (192.168.X.X)
				if strings.HasPrefix(ipv4.String(), "192.168.") {
					return ipv4.String()
				}
			}
		}
	}

	// Si aún no encontramos, tomamos cualquiera que no sea loopback
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipv4 := ipNet.IP.To4(); ipv4 != nil {
				return ipv4.String()
			}
		}
	}

	return "localhost"
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/download/", downloadHandler)
	http.HandleFunc("/upload", uploadHandler)
	// http.HandleFunc("/static/", staticHandler)

	// Crea un sub-filesystem que solo ve el directorio "static"
	staticFS, err := fs.Sub(embeddedFiles, "static")
	if err != nil {
		panic(err)
	}

	// Sirve los archivos estáticos desde el sub-filesystem
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	// http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	localIP := getLocalIP()
	port := "8080"

	fmt.Printf("Serving on:\n")
	fmt.Printf("➡ Local:   http://localhost:%s\n", port)
	fmt.Printf("➡ Network: http://%s:%s\n", localIP, port)

	http.ListenAndServe(":8080", nil)
}

package main

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

//go:embed templates/* static/*
var embeddedFiles embed.FS

type FileInfo struct {
	Name string
	Size string
}

const (
	// maxUploadSize = 500 << 20 // 500 MB
	maxUploadSize = 5 << 30    // 5 GB
	bufferSize    = 32 * 1024  // 32 KB buffer for file operations
)

// formatSize convierte bytes a formato legible (KB, MB, GB)
func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
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
			log.Printf("Warning: failed to get info for %s: %v", file.Name(), err)
			continue
		}
		fileInfoList = append(fileInfoList, FileInfo{
			Name: file.Name(),
			Size: formatSize(info.Size()),
		})
	}
	return fileInfoList, nil
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	files, err := listFiles()
	if err != nil {
		log.Printf("Error listing files: %v", err)
		http.Error(w, "Failed to list files", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFS(embeddedFiles, "templates/index.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Failed to load template", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, files); err != nil {
		log.Printf("Error executing template: %v", err)
	}
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	fileName := strings.TrimPrefix(r.URL.Path, "/download/")
	
	// Validación de seguridad: prevenir directory traversal
	if fileName == "" || strings.Contains(fileName, "..") || strings.Contains(fileName, "/") || strings.Contains(fileName, "\\") {
		http.Error(w, "Invalid file name", http.StatusBadRequest)
		return
	}

	filePath := filepath.Join(".", fileName)

	// Validar que el archivo existe y es un archivo regular
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("Error stating file %s: %v", fileName, err)
		http.Error(w, "Error accessing file", http.StatusInternalServerError)
		return
	}
	if info.IsDir() {
		http.Error(w, "Cannot download directory", http.StatusBadRequest)
		return
	}

	log.Printf("Downloading file: %s (%s)", fileName, formatSize(info.Size()))
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fileName))
	http.ServeFile(w, r, filePath)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Limitar el tamaño del request body
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		log.Printf("Error parsing multipart form: %v", err)
		http.Error(w, "File too large or invalid form data", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		log.Printf("Error reading file: %v", err)
		http.Error(w, "Failed to read file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validar nombre de archivo
	fileName := filepath.Base(header.Filename)
	if fileName == "" || fileName == "." || fileName == ".." || strings.Contains(fileName, "/") || strings.Contains(fileName, "\\") {
		http.Error(w, "Invalid file name", http.StatusBadRequest)
		return
	}

	dst, err := os.Create(fileName)
	if err != nil {
		log.Printf("Error creating file %s: %v", fileName, err)
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Usar buffer optimizado para la copia
	written, err := io.CopyBuffer(dst, file, make([]byte, bufferSize))
	if err != nil {
		log.Printf("Error saving file %s: %v", fileName, err)
		os.Remove(fileName) // Limpiar archivo parcial
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	log.Printf("File uploaded: %s (%s)", fileName, formatSize(written))
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	fileName := r.FormValue("filename")
	
	// Validación de seguridad
	if fileName == "" || strings.Contains(fileName, "..") || strings.Contains(fileName, "/") || strings.Contains(fileName, "\\") {
		http.Error(w, "Invalid file name", http.StatusBadRequest)
		return
	}

	filePath := filepath.Join(".", fileName)

	// Validar que el archivo existe
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("Error stating file %s: %v", fileName, err)
		http.Error(w, "Error accessing file", http.StatusInternalServerError)
		return
	}
	if info.IsDir() {
		http.Error(w, "Cannot delete directory", http.StatusBadRequest)
		return
	}

	if err := os.Remove(filePath); err != nil {
		log.Printf("Error deleting file %s: %v", fileName, err)
		http.Error(w, "Failed to delete file", http.StatusInternalServerError)
		return
	}

	log.Printf("File deleted: %s", fileName)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// getLocalIP obtiene la IP local preferida para la red
func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Printf("Error getting network interfaces: %v", err)
		return "localhost"
	}

	// Preferencia: 192.168.x.x (redes domésticas típicas)
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipv4 := ipNet.IP.To4(); ipv4 != nil {
				if strings.HasPrefix(ipv4.String(), "192.168.") {
					return ipv4.String()
				}
			}
		}
	}

	// Segunda opción: cualquier IPv4 privada que no sea virtual
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipv4 := ipNet.IP.To4(); ipv4 != nil {
				ip := ipv4.String()
				// Excluir redes virtuales comunes
				if !strings.HasPrefix(ip, "172.17.") && // Docker
					!strings.HasPrefix(ip, "192.168.56.") { // VirtualBox
					return ip
				}
			}
		}
	}

	// Última opción: cualquier IPv4 no loopback
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
	// Configurar puerto (puede ser sobrescrito con variable de entorno PORT)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Configurar rutas
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/download/", downloadHandler)
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/delete", deleteHandler)

	// Configurar archivos estáticos embebidos
	staticFS, err := fs.Sub(embeddedFiles, "static")
	if err != nil {
		log.Fatalf("Error loading static files: %v", err)
	}
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	// Crear servidor con timeouts
	srv := &http.Server{
		Addr:         ":" + port,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Mostrar información de conexión
	localIP := getLocalIP()
	fmt.Println("\n🚀 File Sharing Server Started")
	fmt.Println("═══════════════════════════════")
	fmt.Printf("➡  Local:   http://localhost:%s\n", port)
	fmt.Printf("➡  Network: http://%s:%s\n", localIP, port)
	fmt.Println("═══════════════════════════════")
	fmt.Printf("📁 Serving files from: %s\n", getCurrentDir())
	fmt.Printf("📤 Max upload size: %s\n", formatSize(maxUploadSize))
	fmt.Println("\nPress Ctrl+C to stop\n")

	// Configurar graceful shutdown
	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		fmt.Println("\n\n🛑 Shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("Server shutdown error: %v", err)
		}
		close(idleConnsClosed)
	}()

	log.Printf("Server starting on port %s", port)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}

	<-idleConnsClosed
	fmt.Println("✅ Server stopped gracefully")
}

// getCurrentDir obtiene el directorio actual de trabajo
func getCurrentDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return "."
	}
	return dir
}

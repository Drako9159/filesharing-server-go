# 🚀 File Sharing Server (Go)

Un servidor de compartición de archivos simple, rápido y seguro escrito en Go. Perfecto para compartir archivos rápidamente en tu red local.

## ✨ Características

- 📂 **Listar archivos** del directorio actual
- ⬇️ **Descargar archivos** con un solo clic
- 📤 **Subir archivos** (hasta 500 MB)
- 🗑️ **Eliminar archivos** con confirmación
- 🎨 **Interfaz moderna** y responsive
- 🔒 **Seguridad mejorada** contra directory traversal
- 🌐 **Acceso en red local** automático
- 💾 **Archivos estáticos embebidos** (no requiere archivos externos)
- ⚡ **Graceful shutdown** (cierre ordenado con Ctrl+C)
- 📊 **Formato de tamaños inteligente** (B, KB, MB, GB)

## 🔒 Mejoras de Seguridad Implementadas

- ✅ Validación de rutas para prevenir directory traversal
- ✅ Límite de tamaño de upload (500 MB)
- ✅ Validación estricta de nombres de archivo
- ✅ Headers de seguridad apropiados
- ✅ Timeouts configurados en el servidor (15s read/write, 60s idle)
- ✅ Cleanup automático de archivos parciales en caso de error

## 🚀 Uso Rápido

### 1. Ejecutar directamente con Go

```bash
go run main.go
```

### 2. Compilar y ejecutar

```bash
go build -o filesharing-server
./filesharing-server
```

El servidor iniciará en el puerto 8080 por defecto y mostrará:
```
🚀 File Sharing Server Started
═══════════════════════════════
➡  Local:   http://localhost:8080
➡  Network: http://192.168.1.X:8080
═══════════════════════════════
📁 Serving files from: /current/directory
📤 Max upload size: 500.00 MB

Press Ctrl+C to stop
```

### 3. Puerto personalizado

```bash
PORT=3000 ./filesharing-server
```

## 📦 Compilación para Windows

### Generar archivo de recursos (opcional - para icono)

```bash
$(go env GOPATH)/bin/rsrc -ico home_server.ico -manifest app.manifest -o rsrc.syso
```

### Compilar para Windows desde Linux/Mac

```bash
GOOS=windows GOARCH=amd64 go build -o filesharing-server.exe
```

### Compilar para Windows en Windows

```bash
go build -o filesharing-server.exe
```

## 📝 Endpoints API

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| GET | `/` | Página principal con listado de archivos |
| GET | `/download/{filename}` | Descargar un archivo específico |
| POST | `/upload` | Subir archivo (multipart/form-data) |
| POST | `/delete` | Eliminar archivo (form: filename) |
| GET | `/static/*` | Archivos estáticos (CSS, JS, imágenes) |

## 🛠️ Optimizaciones Implementadas

### 1. Performance
- ⚡ Buffer optimizado de 32 KB para transferencias de archivos
- ⚡ Timeouts configurados para prevenir conexiones colgadas
- ⚡ Uso eficiente de `io.CopyBuffer` para streaming

### 2. Código Limpio
- ✨ Eliminación de código comentado y funciones no usadas
- ✨ Funciones mejor organizadas y documentadas
- ✨ Logging estructurado con contexto
- ✨ Constantes para valores mágicos

### 3. Experiencia de Usuario
- 🎨 UI moderna con degradados y animaciones
- 📱 Diseño responsive (móvil, tablet, desktop)
- 😊 Emojis informativos en la consola
- ⚠️ Confirmación antes de eliminar archivos
- 📊 Formato de tamaños automático y legible

### 4. Manejo de Errores
- 🐛 Logging detallado de todos los errores
- 💬 Mensajes de error claros para el usuario
- 🧹 Cleanup automático de archivos parciales
- 🛡️ Validación exhaustiva de inputs

## 🎯 Casos de Uso

- 🏠 Compartir archivos entre dispositivos en tu casa
- 💼 Transferencias rápidas en oficina/LAN
- 🧪 Testing y desarrollo de aplicaciones
- 📱 Enviar archivos a móviles sin cables
- 🎮 Compartir mods, saves, o archivos de juegos

## 🔧 Requisitos

- Go 1.20 o superior
- Puerto 8080 disponible (o configurable con variable PORT)

### Instalación de Go (si no lo tienes)

```bash
# Linux/Mac
wget https://go.dev/dl/go1.24.1.linux-amd64.tar.gz
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.24.1.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
```

## 📋 Instalación

```bash
# Clonar el repositorio
git clone https://github.com/Drako9159/filesharing-server-go.git
cd filesharing-server-go

# Ejecutar directamente
go run main.go

# O compilar primero
go build -o filesharing-server
./filesharing-server
```

## ⚠️ Nota de Seguridad

Este servidor está diseñado para uso en **redes locales confiables**. 

**NO expongas este servidor directamente a Internet** sin implementar:
- 🔐 Autenticación de usuarios
- 🔒 HTTPS/TLS
- 🧱 Firewall y rate limiting
- 🛡️ Protección adicional contra ataques

## 🤝 Contribuir

Las contribuciones son bienvenidas. Por favor:

1. Fork el proyecto
2. Crea una rama para tu feature (`git checkout -b feature/AmazingFeature`)
3. Commit tus cambios (`git commit -m 'Add some AmazingFeature'`)
4. Push a la rama (`git push origin feature/AmazingFeature`)
5. Abre un Pull Request

## 📄 Licencia

Este proyecto es de código abierto y está disponible para uso personal y educativo.

## 🙏 Agradecimientos

Desarrollado con ❤️ usando Go y HTML/CSS moderno.
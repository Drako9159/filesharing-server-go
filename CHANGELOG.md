# Changelog - v2.0.0

All notable changes to the **My Local Cloud (Go)** project are documented in this file.
## [2.0.0] - 2026-04-21
### 🚀 Added**Security Layer**: Integrated a password-protected login system using session-based cookies (`auth_token`).
**Dark Mode UI**: Completely redesigned the frontend with a modern, eye-friendly Dark Mode aesthetic.
**Real-time Search**: Added a search bar that filters files instantly as the user types.
**Progress Tracking**: New visual progress bar for file uploads with real-time percentage feedback.
**Safe Deletion**: Integrated file removal functionality with a confirmation popup to prevent accidental data loss.
**System Protection**: Implemented an exclusion list to hide sensitive files (`main.go`, `go.mod`, `.exe`, etc.) from the public web interface.

### ⚡ Performance & Optimization
- **Zero-RAM Streaming**: Re-engineered the upload engine to stream data directly to disk. This allows handling massive files (10GB+) even on hardware
      with very limited RAM.
- **Download Stability**: Removed fixed `WriteTimeouts` to support stable downloads of large files over WiFi and slower network connections, fixing the
      previous "keep file" browser prompts.
- **Mobile-First Design**: Replaced legacy HTML tables with a CSS Grid/Card layout for perfect responsiveness on smartphones and tablets.

### 📝 Documentation
- **Global Reach**: Migrated the primary documentation (`README.md`) to English.
- **Bilingual Support**: Preserved the original Spanish documentation as `README_SPANISH.md`.
### 🔧 Build System
- **Cross-Compilation**: Optimized build commands for generating Windows `.exe` binaries from Linux environments.
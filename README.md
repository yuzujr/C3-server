# C3 Server
C3 Server is a Go-based backend for the C3 framework, offering better scalability, easier deployment, and higher reliability compared to previous versions.

## Work in progress
This project is currently under development.

### ✅ Completed

* Basic client upload API implemented
* Implement WebSocket communication with clients

### 🔜 Upcoming

* Connect and sync with web frontend

### 🔭 Outlook

The current frontend code is retained for now, but it has several limitations. In the future, I plan to refactor it using modern frameworks such as React, Vue, or Svelte. For now, our main focus is on backend development, and the frontend overhaul will be addressed in later stages.

## Features
- **WebSocket Support**: Real-time communication with clients.
- **Web frontend**: A simple web interface for managing clients.
- **Database Integration**: Uses PostgreSQL for persistent storage.
- **Environment Configuration**: Uses `.env` files for easy configuration management.

## Quick Start
1. **Clone the repository**
   ```sh
   git clone https://github.com/yuzujr/C3-server.git
   cd C3-server
   ```

2. **Configure environment variables**
   - Copy `.env.example` to `.env` and update database and other settings as needed.

3. **Start the server**
   ```sh
   # If you are unable to connect to the official Go module proxy, you can set an alternative proxy (e.g., `https://goproxy.io`) to download dependencies:
   # go env -w GOPROXY=https://goproxy.io,direct
   go mod tidy
   go run cmd/server/main.go
   ```
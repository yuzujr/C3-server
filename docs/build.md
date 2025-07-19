## Server build
1. **Clone the repository**

   ```sh
   git clone https://github.com/yuzujr/C3-server.git
   cd C3-server
   ```

2. **build**

   ```sh
   # If you are unable to connect to the official Go module proxy,
   # you can set an alternative proxy (e.g., `https://goproxy.io`) to download dependencies:
   # go env -w GOPROXY=https://goproxy.io,direct
   go mod tidy
   # Linux:
   go run ./cmd/server/
   # Windows:
   go run .\cmd\server\
   ```
# Development mode - run without building frontend
dev:
	go run main.go

# Build frontend assets
build-frontend:
	@echo "Building frontend..."
	cd dashboard && npm run build
	@echo "Frontend build complete → server/dist"

# Generate Windows resources
gen:
	go generate

# Full build for Windows (frontend + backend)
build: build-frontend
	@echo "Building Windows executable..."
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc go build -ldflags="-s -w" -o antower.exe main.go
	@echo "Windows build complete → antower.exe"

# Clean build artifacts
clean:
	rm -rf server/dist
	rm -f antower.exe

.PHONY: dev build-frontend gen build-win build clean

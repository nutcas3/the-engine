BINARY_NAME=engine-fn
CLI_NAME=engine-cli
LDFLAGS=-s -w -X main.version=$(shell git describe --tags --always --dirty)
GCFLAGS=all=-l
TRIMPATH=-trimpath
BUILDID=-buildid=

.PHONY: build build-cli build-web encrypt clean test dev install-deps docker-up docker-caddy docker-traefik docker-nginx build-minimal build-nocompress

build:
	@echo "Building Sovereign Engine..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="$(LDFLAGS) $(BUILDID)" -gcflags="$(GCFLAGS)" $(TRIMPATH) -o bin/$(BINARY_NAME) ./cmd/function
	@echo "Compressing binary..."
	upx --best --lzma bin/$(BINARY_NAME)

build-cli:
	@echo "Building Sovereign Engine CLI..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="$(LDFLAGS) $(BUILDID)" -gcflags="$(GCFLAGS)" $(TRIMPATH) -o bin/$(CLI_NAME) ./cmd/cli
	@echo "Compressing binary..."
	upx --best --lzma bin/$(CLI_NAME)

build-web:
	@echo "Building Sovereign Engine Web Server..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="$(LDFLAGS) $(BUILDID)" -gcflags="$(GCFLAGS)" $(TRIMPATH) -o bin/engine-web ./cmd/ui
	@echo "Compressing binary..."
	upx --best --lzma bin/engine-web

build-minimal:
	@echo "Building Sovereign Engine (minimal size)..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="$(LDFLAGS) $(BUILDID)" -gcflags="$(GCFLAGS)" $(TRIMPATH) -o bin/$(BINARY_NAME) ./cmd/function
	@echo "Compressing binary with maximum compression..."
	upx --ultra-brute bin/$(BINARY_NAME)

build-nocompress:
	@echo "Building Sovereign Engine (no compression)..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="$(LDFLAGS) $(BUILDID)" -gcflags="$(GCFLAGS)" $(TRIMPATH) -o bin/$(BINARY_NAME) ./cmd/function
	@echo "Binary built without UPX compression"

encrypt:
	@echo "Building encryption tool..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="$(LDFLAGS) $(BUILDID)" -gcflags="$(GCFLAGS)" $(TRIMPATH) -o bin/encrypt ./cmd/encrypt
	@echo "Compressing binary..."
	upx --best --lzma bin/encrypt
	@echo "Encryption tool built as bin/encrypt"

clean:
	rm -rf bin/

test:
	go test ./...

dev:
	go run ./cmd/function

install-deps:
	go mod tidy
	go mod download

deploy:
	kubectl apply -f apis/
	kubectl apply -f compositions/
	kubectl apply -f configs/

docker-up:
	docker-compose up -d

docker-caddy:
	WEBSERVER=caddy docker-compose up -d

docker-traefik:
	WEBSERVER=traefik docker-compose up -d

docker-nginx:
	WEBSERVER=nginx docker-compose up -d

docker-down:
	docker-compose down

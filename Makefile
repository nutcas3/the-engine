BINARY_NAME=engine-fn
LDFLAGS=-s -w

.PHONY: build clean test dev install-deps docker-up docker-caddy docker-traefik docker-nginx

build:
	@echo "Building Sovereign Engine..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o bin/$(BINARY_NAME) ./cmd/function
	@echo "Compressing binary..."
	upx --brute bin/$(BINARY_NAME)

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

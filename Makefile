.PHONY: serve
serve:
	go run main.go serve

.PHONY: build-linux-amd64
build-linux-amd64:
	go version
	go mod tidy
	GOARCH=amd64 GOOS=linux go build -o build/lite-ai-gateway-linux-amd64 -ldflags "-s -w" main.go

.PHONY: build-linux-arm64
build-linux-arm64:
	go version
	go mod tidy
	GOARCH=arm64 GOOS=linux go build -o build/lite-ai-gateway-linux-arm64 -ldflags "-s -w" main.go
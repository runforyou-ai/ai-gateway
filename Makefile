.PHONY: generate serve build-linux-amd64 build-linux-arm64 tidy

generate:
	go generate ./internal/contract

serve:
	go run . serve

tidy:
	go mod tidy

build-linux-amd64: generate
	go version
	go mod tidy
	GOARCH=amd64 GOOS=linux go build -o build/lite-ai-gateway-linux-amd64 -ldflags "-s -w" .

build-linux-arm64: generate
	go version
	go mod tidy
	GOARCH=arm64 GOOS=linux go build -o build/lite-ai-gateway-linux-arm64 -ldflags "-s -w" .

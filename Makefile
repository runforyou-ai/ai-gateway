.PHONY: dev build generate amd64 arm64

dev:
	go run . serve &
	cd frontend && pnpm dev

generate:
	go generate ./internal/contract

build:
	@arch="$(filter amd64 arm64,$(MAKECMDGOALS))"; \
	set -e; \
	$(MAKE) generate; \
	go mod tidy; \
	( cd frontend && pnpm build ); \
	rm -rf internal/web/dist; \
	cp -R frontend/dist internal/web/dist; \
	mkdir -p build; \
	if [ -z "$$arch" ]; then \
		go build -o build/lite-ai-gateway -ldflags "-s -w" .; \
	elif [ "$$(echo $$arch | wc -w | tr -d " ")" -ne 1 ]; then \
		echo "usage: make build [amd64|arm64]"; \
		exit 1; \
	else \
		GOARCH=$$arch GOOS=linux go build -o build/lite-ai-gateway-linux-$$arch -ldflags "-s -w" .; \
	fi

amd64 arm64:
	@:

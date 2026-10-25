.PHONY: dev backend frontend build amd64 arm64

dev:
	@target="$(filter backend frontend,$(MAKECMDGOALS))"; \
	if [ "$$target" = "backend" ]; then \
		go generate ./internal/contract && go run . serve; \
	elif [ "$$target" = "frontend" ]; then \
		cd frontend && pnpm dev; \
	else \
		echo "usage: make dev backend|frontend"; \
		exit 1; \
	fi

backend frontend:
	@:

build:
	@arch="$(filter amd64 arm64,$(MAKECMDGOALS))"; \
	set -e; \
	go generate ./internal/contract; \
	go mod tidy; \
	( cd frontend && pnpm build ); \
	rm -rf internal/web/dist; \
	cp -R frontend/dist internal/web/dist; \
	mkdir -p build; \
	if [ -z "$$arch" ]; then \
		go build -o build/ai-gateway -ldflags "-s -w" .; \
	elif [ "$$(echo $$arch | wc -w | tr -d " ")" -ne 1 ]; then \
		echo "usage: make build [amd64|arm64]"; \
		exit 1; \
	else \
		GOARCH=$$arch GOOS=linux go build -o build/ai-gateway-linux-$$arch -ldflags "-s -w" .; \
	fi

amd64 arm64:
	@:

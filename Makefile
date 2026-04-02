# ============================================================
#  SliceFlow — Makefile
# ============================================================

APP_NAME  = sliceflow
CMD_PATH  = ./cmd/api
BUILD_DIR = ./bin

.PHONY: run watch build test test-cover lint fmt vet swagger docker-up docker-down docker-logs clean

## ── Dev ─────────────────────────────────────────────────────

run:
	go run $(CMD_PATH)/main.go

# Live-reload: go install github.com/air-verse/air@latest
watch:
	air

## ── Build ───────────────────────────────────────────────────

build:
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) $(CMD_PATH)

## ── Testing ─────────────────────────────────────────────────

test:
	go test ./... -race

test-cover:
	go test ./... -race -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

## ── Code Quality ────────────────────────────────────────────

# go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
lint:
	golangci-lint run ./...

fmt:
	gofmt -w .

vet:
	go vet ./...

## ── Swagger ─────────────────────────────────────────────────

# go install github.com/swaggo/swag/cmd/swag@latest
swagger:
	swag init -g $(CMD_PATH)/main.go --output ./Sliceflow/docs

## ── Docker ──────────────────────────────────────────────────

docker-up:
	docker compose up -d --build

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f

## ── Cleanup ─────────────────────────────────────────────────

clean:
	rm -rf $(BUILD_DIR) coverage.out coverage.html

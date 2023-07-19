build:
	go build -o bin/guess cmd/game/main.go

fmt:
	gofumpt -w .

tidy:	
	go mod tidy

lint: build fmt tidy
	golangci-lint run ./...

run:
	go run cmd/game/main.go
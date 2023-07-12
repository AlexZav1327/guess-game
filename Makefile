build:
	go build cmd/game/main.go

fmt:
	gofumpt -w .

tidy:	
	go mod tidy

lint: 
	golangci-lint run ./...
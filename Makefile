build:
	go build cmd/cosmos-transactions-bot.go

install:
	go install cmd/cosmos-transactions-bot.go

lint:
	golangci-lint run --fix ./...

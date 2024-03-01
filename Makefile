build:
	go build -o bin/cron cmd/cron/main.go

run:
	go run cmd/cron/main.go

test:
	go test -v ./...
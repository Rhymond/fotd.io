get:
	go get ./...

build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/get get/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/fetch fetch/main.go
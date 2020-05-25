.PHONY: build
build:
	go build -o ./bin/server ./cmd/server/main.go
	go build -o ./bin/news ./cmd/news/main.go
	go build -o ./bin/auth ./cmd/auth/main.go
	go build -o ./bin/user ./cmd/user/main.go
	go build -o ./bin/gitserver ./cmd/gitserver/main.go

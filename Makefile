.PHONY: build
build:
	go build -o ./bin/server ./cmd/app/server.go
	go build -o ./bin/admin ./cmd/admin_service/admin_main.go
	go build -o ./bin/order ./cmd/order_main/order_main.go
	go build -o ./bin/session ./cmd/session_service/session_main.go

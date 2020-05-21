.PHONY: build
build:
	go build -o ./bin/server ./cmd/app/server.go
	go build -o ./bin/admin ./cmd/admin_service/admin_main.go
	go build -o ./bin/order ./cmd/order_main/order_main.go
	go build -o ./bin/session ./cmd/session_service/session_main.go

.PHONY: test
test:
	go test ./...

.PHONY: cover
cover:
	go test -coverpkg=./... ./... -coverprofile cover.out.tmp
	fgrep cover.out.tmp  -e '/mocks/' -v > cover.out.tmp2
	fgrep cover.out.tmp2 -e '/docs/' -v > cover.out.tmp3
	fgrep cover.out.tmp3 -e '_proto.pb.go' -v > cover.out.tmp4
	fgrep cover.out.tmp4 -e '_easyjson.go' -v > cover.out
	go tool cover -func cover.out | grep total
	rm -f cover.out.tmp*
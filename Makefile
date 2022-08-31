.PHONY: test show-coverage run file-run db-run run-ld proto db-with-grpc

test:
	 go test ./... -coverprofile cp.out

show-coverage:
	go tool cover -func cp.out | grep total:

run:
	go run cmd/shortener/main.go

file-run:
	go run cmd/shortener/main.go -f="./tmp"

db-run:
	go run cmd/shortener/main.go -d="postgres://root:root@localhost:5432/postgres?sslmode=disable"

run-ld:
	go run -ldflags "-X main.buildVersion=v1.0.0 -X 'main.buildDate=$(date +'%Y/%m/%d %H:%M:%S')' -X 'main.buildCommit='" cmd/shortener/main.go

run-checks:
	go run cmd/staticlint/main.go -builtin -static -extra ./...

proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative internal/app/grpc/proto/api.proto

db-with-grpc:
	go run cmd/shortener/main.go -d="postgres://root:root@localhost:5432/postgres?sslmode=disable" -g="3200"
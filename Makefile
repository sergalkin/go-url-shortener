.PHONY: test show-coverage run file-run db-run

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
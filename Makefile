run:
	go run .

build:
	go build -o go-cat .

test:
	go test cat_test.go

format:
	gofmt -w ./*.go

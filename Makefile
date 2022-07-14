run:
	go run main.go

linter:
	golangci-lint run

all: linter run
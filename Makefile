run:
	go run main.go

linter:
	golangci-lint run

swagger:
	GO111MODULE=off swagger generate spec -o ./swagger.yaml --scan-models

all: linter run
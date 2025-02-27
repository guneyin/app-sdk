init: clean tidy vet

clean:
	go clean

tidy:
	go mod tidy

vet:
	go vet ./...

lint:
	golangci-lint run

fix:
	golangci-lint run --fix

protoc:
	protoc ./proto/*.proto --go_out=./proto --go-grpc_out=./proto

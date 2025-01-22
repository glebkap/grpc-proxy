
all: gen

gen: pb/test.proto
	cd pb && protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    test.proto

run: 
	go run cmd/main.go

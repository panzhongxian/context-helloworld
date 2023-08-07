all: build

build:
	go mod tidy
	go build -o server greeter_server/main.go
	go build -o client greeter_client/main.go
	go build -o server_with_otel greeter_server_with_otel/main.go
	go build -o client_with_otel greeter_client_with_otel/main.go

clean:
	@ rm server client server_with_otel client_with_otel


.PHONY: build clean 

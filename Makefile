all: clean build

build:
	go build -o server greeter_server/main.go
	go build -o client greeter_client/main.go

clean:
	@ rm server client


.PHONY: build clean 

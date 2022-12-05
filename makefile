run:
	go run main.go

build:
	go build
	joybox.exe

test:
	go test ./... -cover
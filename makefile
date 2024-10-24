compile:
	GOOS=darwin GOARCH=amd64 go build -o ./build/app_amd64 ./cmd/main.go

run: compile
	./build/app_amd64
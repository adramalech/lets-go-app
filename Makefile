#######################################################################################################################
#
# makefile docs:  https://tutorialedge.net/golang/makefiles-for-go-developers/
#
# GOARCH & GOOS:  https://gist.github.com/asukakenji/f15ba7e588ac42795f421b48b8aede63
#
#######################################################################################################################

BUILD_PATH=./cmd/web/*

build:
	@echo "Building app..."
	go build -o bin/snippetbox $(BUILD_PATH)

run:
	@echo "Run locally..."
	go run ./cmd/web/* -addr=":10000" --static-dir="./ui/static"

prod:
	@echo "Run production release..."
	go run ./bin/snippetbox --addr=":10000" --static-dir="./ui/static" >> ./tmp/info.log 2 >> ./tmp/error.log

compile:
	@echo "Cross compile..."
	GOOS=linux GOARCH=amd64 go build -o bin/snippetbox-linux-amd64 $(BUILD_PATH)
	GOOS=windows GOARCH=amd64 go build -o bin/snippetbox-windows-amd64.exe $(BUILD_PATH)
	GOOS=darwin GOARCH=amd64 go build -o bin/snippetbox-darwin-amd64 $(BUILD_PATH)
	GOOS=linux GOARCH=arm64 go build -o bin/snippetbox-linux-arm64 $(BUILD_PATH)

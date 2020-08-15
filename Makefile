#######################################################################################################################
#
# makefile docs:  https://tutorialedge.net/golang/makefiles-for-go-developers/
#
# GOARCH & GOOS:  https://gist.github.com/asukakenji/f15ba7e588ac42795f421b48b8aede63
#
#######################################################################################################################

build:
	go build -o bin/snippetbox ./cmd/web* 

run:
	go run ./cmd/web/* -addr=":10000" --static-dir="./ui/static"

prod:
	go run ./bin/snippetbox --addr=":10000" --static-dir="./ui/static" >> ./tmp/info.log 2>>./tmp/error.log

compile:
	GOOS=freebsd GOARCH=386 go build -o bin/snippetbox-freebsd-386 ./cmd/web/*
	GOOS=linux GOARCH=386 go build -o bin/snippetbox-linux-386 ./cmd/web/*
	GOOS=windows GOARCH=386 go build -o bin/snippetbox-windows-386.exe ./cmd/web/*
	GOOS=darwin GOARCH=386 go build -o bin/snippetbox-darwin-386 ./cmd/web/*
	GOOS=freebsd GOARCH=amd64 go build -o bin/snippetbox-freebsd-amd64 ./cmd/web/*
	GOOS=linux GOARCH=amd64 go build -o bin/snippetbox-linux-amd64 ./cmd/web/*
	GOOS=windows GOARCH=amd64 go build -o bin/snippetbox-windows-amd64.exe ./cmd/web/*
	GOOS=darwin GOARCH=amd64 go build -o bin/snippetbox-darwin-amd64 ./cmd/web/*
	GOOS=freebsd GOARCH=arm go build -o bin/snippetbox-freebsd-arm ./cmd/web/*
	GOOS=linux GOARCH=arm go build -o bin/snippetbox-linux-arm ./cmd/web/*
	GOOS=windows GOARCH=arm go build -o bin/snippetbox-windows-arm.exe ./cmd/web/*
	GOOS=freebsd GOARCH=arm64 go build -o bin/snippetbox-freebsd-arm64 ./cmd/web/*
	GOOS=linux GOARCH=arm64 go build -o bin/snippetbox-linux-arm64 ./cmd/web/*

#######################################################################################################################
#
# makefile docs:  https://tutorialedge.net/golang/makefiles-for-go-developers/
#
# GOARCH & GOOS:  https://gist.github.com/asukakenji/f15ba7e588ac42795f421b48b8aede63
#
#######################################################################################################################

BUILD_PATH=./cmd/web/*
DB_CONN="web:password12345@127.0.0.1:8080/snippetbox?parseTime=true"
TAG_VERSION=0.0.2

build:
	go build -o ./bin/snippetbox $(BUILD_PATH)

docker-build:
	docker build -t snippetbox-go .

docker-publish:
	docker tag snippetbox-go jthrone/snippetbox-go:$(TAG_VERSION)
	docker login
	docker push jthrone/snippetbox-go:$(TAG_VERSION)

docker: docker-build docker-publish

run:
	@echo "Run locally..."
	go run ./cmd/web/* -addr=":10000" -dsn=$(DB_CONN) --static-dir="./ui/static"

deploy:
	@echo "Setup kubernetes cluster locally for minikube..."
	@echo ""
	@echo "Apply secrets..."
	kubectl apply -f ./mysql-secret.yaml
	@echo ""
	@echo "Apply mysql database pod..."
	kubectl apply -f ./mysql.yaml
	@echo ""
	@echo "Build app docker container..."
	docker build -t snippetbox-go
	@echo ""
	@echo "Apply app pod..."
	kubectl apply -f ./snippetbox.yaml

prod:
	go run ./bin/snippetbox -addr=":$(APP_PORT)" -dsn="$(MYSQL_USERNAME):$(MYSQL_PASSWORD)@$(MYSQL_DATABASE_HOST):$(MYSQL_DATABASE_PORT)/$(MYSQL_DATABASE_NAME)?parseTime=true" --static-dir="./ui/static"

compile:
	@echo "Cross compile..."
	GOOS=linux GOARCH=amd64 go build -o bin/snippetbox-linux-amd64 $(BUILD_PATH)
	GOOS=windows GOARCH=amd64 go build -o bin/snippetbox-windows-amd64.exe $(BUILD_PATH)
	GOOS=darwin GOARCH=amd64 go build -o bin/snippetbox-darwin-amd64 $(BUILD_PATH)
	GOOS=linux GOARCH=arm64 go build -o bin/snippetbox-linux-arm64 $(BUILD_PATH)

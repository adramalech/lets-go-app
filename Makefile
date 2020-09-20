#######################################################################################################################
#
# makefile docs:  https://tutorialedge.net/golang/makefiles-for-go-developers/
#
# GOARCH & GOOS:  https://gist.github.com/asukakenji/f15ba7e588ac42795f421b48b8aede63
#
#######################################################################################################################

BUILD_PATH=./cmd/web/*
DB_CONN="web:password12345@127.0.0.1:8080/snippetbox?parseTime=true"
TAG_VERSION=0.0.6

build:
	go build -o ./bin/snippetbox $(BUILD_PATH)

docker-build:
	docker build -t snippetbox-go .

docker-publish:
	docker tag snippetbox-go jthrone/snippetbox-go:$(TAG_VERSION)
	docker login
	docker push jthrone/snippetbox-go:$(TAG_VERSION)

docker: docker-build docker-publish

# docker-private-repo-secret:
#	docker login
#	kubectl create secret generic regcred --from-file=.dockerconfigjson=/Users/jonathanthrone/.docker/config.json --type=kubernetes.io/dockerconfigjson
#	kubectl get secret regcred --output=yaml

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
	@echo "Apply mysql configmap..."
	kubectl apply -f ./mysql-configmap.yaml
	@echo ""
	@echo "Apply app pod..."
	kubectl apply -f ./snippetbox.yaml

prod:
	./bin/snippetbox -addr=":$(APP_PORT)" -dsn="$(MYSQL_USERNAME):$(MYSQL_PASSWORD)@$(MYSQL_DATABASE_HOST):$(MYSQL_DATABASE_PORT)/$(MYSQL_DATABASE_NAME)?parseTime=true" --static-dir="./ui/static"

destroy-everything:
	kubectl delete all --all --all-namespaces

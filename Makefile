include .env

##### Arguments ######
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
GOPATH  ?= $(shell go env GOPATH)
ENV_FILE = ./.env
PROJECT_NAME = medioa
COLOR := "\e[1;36m%s\e[0m\n"
DOMAIN = https://medioa.fly.dev


##### Build Flags ######
CGO_ENABLED=0
BUILD_TAG_FLAG=
SERVER_PATH = ./.bin/medioa

##### Run server #####
DOCKER_COMPOSE_FILES := -f ./build/docker-compose.yml

dev:
	go run cmd/main.go

lint:
	go mod tidy
	golangci-lint run
	govulncheck ./...

server:
	@make server-build
	@printf $(COLOR) "Starting server"
	@$(SERVER_PATH)

server-build:
	@printf $(COLOR) "Build $(SERVER_PATH) with CGO_ENABLED=$(CGO_ENABLED) for $(GOOS)/$(GOARCH)..."
	CGO_ENABLED=$(CGO_ENABLED) go build $(BUILD_TAG_FLAG) -o $(SERVER_PATH) ./cmd

doc:
	@printf $(COLOR) "Starting swagger generating"
	swag fmt
	swag init -g cmd/main.go --pd
	@printf $(COLOR) "Swagger generated"

tag:
	@printf $(COLOR) "Release version $(VERSION)"
	git tag -a v$(VERSION) -m "Release version $(VERSION)"
	git push origin v$(VERSION)

docker-tag:
	@printf $(COLOR) "Release image version $(PROJECT_NAME):$(VERSION)"
	docker tag $(PROJECT_NAME):$(VERSION) $(USERNAME)/$(PROJECT_NAME):$(VERSION)
	docker push $(USERNAME)/$(PROJECT_NAME):$(VERSION)

docker-build:
	@printf $(COLOR) "Building docker image $(PROJECT_NAME)"
	docker build -t $(PROJECT_NAME) .

docker-buildx:
	@printf $(COLOR) "Building docker image $(USERNAME)/$(PROJECT_NAME):$(VERSION)"
	docker buildx build --platform linux/amd64 -t $(USERNAME)/$(PROJECT_NAME):$(VERSION) --load .

up:
	docker compose -p $(PROJECT_NAME) $(DOCKER_COMPOSE_FILES) --env-file $(ENV_FILE) up -d
	@printf $(COLOR) "Server started: http://localhost:$(PORT)"
	@printf $(COLOR) "Go to homepage: http://localhost:$(PORT)/index"
	@printf $(COLOR) "Go to swagger: http://localhost:$(PORT)/swagger/index.html"

down:
	docker compose -p $(PROJECT_NAME) $(DOCKER_COMPOSE_FILES) --env-file $(ENV_FILE) down
	@printf $(COLOR) "Server stopped and remove all container"

stop:
	docker compose -p $(PROJECT_NAME) $(DOCKER_COMPOSE_FILES) --env-file $(ENV_FILE) stop
	@printf $(COLOR) "Server stopped"

open:
	@printf $(COLOR) "Opening $(DOMAIN)/index"
	open $(DOMAIN)/index

fly-launch:
	@printf $(COLOR) "Launching $(PROJECT_NAME) to fly.io"
	fly launch

fly-deploy:
	@printf $(COLOR) "Deploying $(PROJECT_NAME) from fly.io"
	fly deploy

fly-secret:
	@printf $(COLOR) "Setting secret $(KEY):$(VALUE) to fly.io"
	fly secrets set $(KEY)=$(VALUE)

fly-secret-del:
	@printf $(COLOR) "Delete secret $(KEY) from fly.io"
	fly secrets unset $(KEY)

fly-secret-view:
	fly secrets list

fly-status:
	fly status

fly-open:
	fly apps open

fly-ips:
	fly ips list

fly-logs:
	fly logs -a $(PROJECT_NAME)
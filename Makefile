include .env

##### Arguments ######
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
GOPATH  ?= $(shell go env GOPATH)
ENV_FILE = ./.env
PROJECT_NAME = medioa
COLOR := "\e[1;36m%s\e[0m\n"


##### Build Flags ######
CGO_ENABLED=0
BUILD_TAG_FLAG=
SERVER_PATH = ./.bin/medioa

##### Run server #####
DOCKER_COMPOSE_FILES := -f ./build/docker-compose.yml

server:
	@make build-server
	@printf $(COLOR) "Starting server"
	@$(SERVER_PATH)

doc:
	@printf $(COLOR) "Starting swagger generating"
	swag fmt
	swag init -g cmd/main.go --pd
	@printf $(COLOR) "Swagger generated"

tag:
	@printf $(COLOR) "Release version $(VERSION)"
	git tag -a v$(VERSION) -m "Release version $(VERSION)"
	git push origin v$(VERSION)

tag-docker:
	@printf $(COLOR) "Release image version $(PROJECT_NAME):v$(VERSION)"
	docker push $(PROJECT_NAME):v$(VERSION)

build-docker:
	@printf $(COLOR) "Building docker image $(PROJECT_NAME)"
	docker build -t $(PROJECT_NAME) .

build-server:
	@printf $(COLOR) "Build $(SERVER_PATH) with CGO_ENABLED=$(CGO_ENABLED) for $(GOOS)/$(GOARCH)..."
	CGO_ENABLED=$(CGO_ENABLED) go build $(BUILD_TAG_FLAG) -o $(SERVER_PATH) ./cmd

up:
	@printf $(COLOR) "Server started: http://localhost:$(APP_PORT)"
	@printf $(COLOR) "Go to homepage: http://localhost:$(APP_PORT)/index"
	@printf $(COLOR) "Go to swagger: http://localhost:$(APP_PORT)/swagger/index.html"
	docker compose -p $(PROJECT_NAME) $(DOCKER_COMPOSE_FILES) --env-file $(ENV_FILE) up -d

down:
	@printf $(COLOR) "Server stopped and remove all container"
	docker compose -p $(PROJECT_NAME) $(DOCKER_COMPOSE_FILES) --env-file $(ENV_FILE) down

stop:
	@printf $(COLOR) "Server stopped"
	docker compose -p $(PROJECT_NAME) $(DOCKER_COMPOSE_FILES) --env-file $(ENV_FILE) stop
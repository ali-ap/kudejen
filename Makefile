# Define variables
PROTOC = protoc
PROTOC_GEN_GO = protoc-gen-go
PROTOC_GEN_GRPC_GO = protoc-gen-go-grpc
GOPATH = $(shell go env GOPATH)
SRC_DIR = ./src/internal/app/proto
OUT_DIR = ./src/internal/app/pb
CMD_DIR = ./src/cmd
BIN_DIR = ./bin
APP_NAME = kudejen
DOCKER = docker
DOCKER_BUILD = $(DOCKER) build
DOCKER_RUN = $(DOCKER) run
DOCKER_IMAGE_NAME = kudejen
DOCKER_CONTAINER_NAME = kudejen-container
# Phony targets
.PHONY: all build clean generate compile run test

# Default target
all: build

# Build target
build: clean generate compile

# Clean target
clean:
	rm -rf $(OUT_DIR)
	mkdir -p $(OUT_DIR)
	rm -rf $(BIN_DIR)

# Generate protobuf target
generate:
	$(PROTOC) --go_out=$(OUT_DIR) --go_opt=paths=source_relative \
	          --go-grpc_out=$(OUT_DIR) --go-grpc_opt=paths=source_relative \
		      --validate_out="lang=go:$(OUT_DIR)" \
	          -I=$(SRC_DIR) $(SRC_DIR)/*.proto

# Compile target
compile:
	go build -o $(BIN_DIR)/$(APP_NAME) $(CMD_DIR)/main.go

# Run target
run:
	go run $(CMD_DIR)/main.go

# Test target
test:
	go test ./...

docker-build:
	$(DOCKER_BUILD) -t $(DOCKER_IMAGE_NAME) .

docker-build-debug:
	$(DOCKER_BUILD) --no-cache -t $(DOCKER_IMAGE_NAME):debug .

docker-run:
	$(DOCKER_RUN) -p 8080:8080 -p 8081:8081 --name $(DOCKER_CONTAINER_NAME) $(DOCKER_IMAGE_NAME)

docker-run-attach:
	$(DOCKER_RUN) -it -p 8080:8080 -p 8081:8081 --name $(DOCKER_CONTAINER_NAME) $(DOCKER_IMAGE_NAME)

docker-stop:
	$(DOCKER) stop $(DOCKER_CONTAINER_NAME)
	$(DOCKER) rm $(DOCKER_CONTAINER_NAME)
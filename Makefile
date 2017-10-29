BUILD_FLAG=--no-cache
PROJECT_NAME=go-warmcache
BUILD_IMAGE_NAME=example-build:latest
BUILD_CANO_IMAGE_NAME="$(PROJECT_NAME)/$(BUILD_IMAGE_NAME)"
RUN_IMAGE_NAME=example:latest
RUN_CANO_IMAGE_NAME="$(PROJECT_NAME)/$(RUN_IMAGE_NAME)"
BINARY_NAME=example

all:install run

install:
	docker build $(BUILD_FLAG) -t $(BUILD_CANO_IMAGE_NAME) -f Dockerfile.build .
	docker run $(BUILD_CANO_IMAGE_NAME) sleep 10 &
	sleep 1
	rm -f ./Dockerfile/$(BINARY_NAME)
	docker cp `docker ps | grep $(BUILD_IMAGE_NAME) | cut -f1 -d' '`\:/go/bin/$(BINARY_NAME) ./Dockerfile/

run:
	docker build $(BUILD_FLAG) -t $(RUN_CANO_IMAGE_NAME) Dockerfile
	docker run $(RUN_CANO_IMAGE_NAME) $(BINARY_NAME)

BUILD_FLAG=--no-cache

USER_NAME=yoheimuta
BUILD_IMAGE_NAME=go-warmcache-example-build:latest
BUILD_CANO_IMAGE_NAME="$(USER_NAME)/$(BUILD_IMAGE_NAME)"

RUN_IMAGE_NAME=go-warmcache-example:latest
RUN_CANO_IMAGE_NAME="$(USER_NAME)/$(RUN_IMAGE_NAME)"
BINARY_NAME=example

CIRCLE_IMAGE_NAME=go-warmcache-circleci:latest
CIRCLE_CANO_IMAGE_NAME="$(USER_NAME)/$(CIRCLE_IMAGE_NAME)"

install:generate run

generate:
	docker build $(BUILD_FLAG) -t $(BUILD_CANO_IMAGE_NAME) -f Dockerfile.build .
	docker run $(BUILD_CANO_IMAGE_NAME) sleep 10 &
	sleep 1
	rm -f ./Dockerfile/$(BINARY_NAME)
	docker cp `docker ps | grep $(BUILD_IMAGE_NAME) | cut -f1 -d' '`\:/go/bin/$(BINARY_NAME) ./Dockerfile/

run:
	docker build $(BUILD_FLAG) -t $(RUN_CANO_IMAGE_NAME) Dockerfile
	docker run $(RUN_CANO_IMAGE_NAME) $(BINARY_NAME)

circle:
	docker build $(BUILD_FLAG) -t $(CIRCLE_CANO_IMAGE_NAME) _circle
	docker push $(CIRCLE_CANO_IMAGE_NAME)

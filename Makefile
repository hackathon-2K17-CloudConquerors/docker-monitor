PROJECT_NAME := docker-monitor
BUILD_NUMBER := test
DOCKER_REGISTRY?=sibiataporeto
DOCKER_IMAGE_NAME?=$(PROJECT_NAME)
DOCKER_IMAGE_TAG?=$(BUILD_NUMBER)

build:
	env GOOS=linux GOARCH=386 go build -o docker-monitor

package: build
	mv docker-monitor docker/docker-monitor

clean:
	rm -rf ./vendor
	rm -rf ./docker/docker-monitor

docker_build: package
		docker \
			build \
			-t $(DOCKER_REGISTRY)/$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG) docker

docker_push: docker_build
		docker \
			push \
			$(DOCKER_REGISTRY)/$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)

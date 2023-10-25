BUILD_NUMBER := latest
PROJECT_NAME := unifi-dns-server
DOCKER_REGISTRY := jodydadescott
DOCKER_IMAGE_NAME?=$(PROJECT_NAME)
DOCKER_IMAGE_TAG?=$(BUILD_NUMBER)

default:
	$(MAKE) container-amd-64

linux-amd-64:
	mkdir -p build/linux-amd-64
	env GOOS=linux GOARCH=arm go build -o build/linux-amd-64/unifi-dns-server unifi-dns-server.go

container-amd-64:
	$(MAKE) linux-amd-64
	docker build -t $(DOCKER_REGISTRY)/$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG) .

push:
	docker push $(DOCKER_REGISTRY)/$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)

clean:
	$(RM) -rf build
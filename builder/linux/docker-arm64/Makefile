BUILD_NUMBER := latest
PROJECT_NAME := home-dns-server-arm64
DOCKER_REGISTRY := jodydadescott
DOCKER_IMAGE_NAME?=$(PROJECT_NAME)
DOCKER_IMAGE_TAG?=$(BUILD_NUMBER)

home-dns-server:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -v -trimpath -o home-dns-server ../../../home-dns-server.go
	docker build -t $(DOCKER_REGISTRY)/$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG) .

push:
	$(MAKE) home-dns-server
	docker push $(DOCKER_REGISTRY)/$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)

clean:
	$(RM) home-dns-server

BUILD_NUMBER := latest
PROJECT_NAME := home-dns-server-amd64
DOCKER_REGISTRY := jodydadescott
DOCKER_IMAGE_NAME?=$(PROJECT_NAME)
DOCKER_IMAGE_TAG?=$(BUILD_NUMBER)

container:
	$(MAKE) home-dns-server
	docker build -t $(DOCKER_REGISTRY)/$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG) .

home-dns-server:
	cd ../bin && $(MAKE)
	cp ../bin/home-dns-server .

push:
	$(MAKE) container
	docker push $(DOCKER_REGISTRY)/$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)

clean:
	$(RM) -f home-dns-server
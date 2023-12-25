all:
	cd amd64 && $(MAKE)
	cd arm64 && $(MAKE)
	cd amd64-docker && $(MAKE)
	cd arm64-docker && $(MAKE)

clean:
	cd amd64 && $(MAKE) clean
	cd arm64 && $(MAKE) clean
	cd amd64-docker && $(MAKE) clean
	cd arm64-docker && $(MAKE) clean

push:
	cd docker-amd64 && $(MAKE) push
	cd docker-arm64 && $(MAKE) push

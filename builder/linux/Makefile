all:
	cd amd64 && $(MAKE)
	cd arm64 && $(MAKE)
	cd docker-amd64 && $(MAKE)
	cd docker-arm64 && $(MAKE)

clean:
	cd amd64 && $(MAKE) clean
	cd arm64 && $(MAKE) clean
	cd docker-amd64 && $(MAKE) clean
	cd docker-arm64 && $(MAKE) clean

push:
	cd docker-amd64 && $(MAKE) push
	cd docker-arm64 && $(MAKE) push

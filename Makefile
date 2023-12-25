default:
	@echo "building for local system; chdir into builder and run make"
	env CGO_ENABLED=0 go build -v -trimpath -o home-dns-server home-dns-server.go
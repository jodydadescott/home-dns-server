default:
	@echo "building for local system; chdir into builder and run make"
	go build -o home-dns-server home-dns-server.go
	
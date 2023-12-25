build:
	env GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -v -trimpath -o home-dns-server ../../../home-dns-server.go

clean:
	$(RM) home-dns-server

test:
	go test -v ./...

clean:
	go clean -modcache

update:
	go get -u ./...

format:
	go fmt ./...
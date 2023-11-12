proto:
	go get github.com/bookpanda/mygraderlist-proto@latest

publish:
	cat ./token.txt | docker login --username bookpanda --password-stdin ghcr.io
	docker build . -t ghcr.io/bookpanda/mygraderlist-auth
	docker push ghcr.io/bookpanda/mygraderlist-auth

test:
	go vet ./...
	go test  -v -coverpkg ./src/app/... -coverprofile coverage.out -covermode count ./src/app/...
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out -o coverage.html

server:
	go run ./src/.

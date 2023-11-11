proto:
	go get github.com/bookpanda/mygraderlist-proto@latest

publish:
	cat ./token.txt | docker login --username bookpanda --password-stdin ghcr.io
	docker build . -t ghcr.io/bookpanda/mygraderlist-auth
	docker push ghcr.io/bookpanda/mygraderlist-auth
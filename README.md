# MyGraderList Auth

MyGraderList is a web app that lets students assess the difficulties and worthiness of each DSA grader problem in their respective courses.

MyGraderList Auth handles the authentication and authorization of the MyGraderList app.

## Technologies

-   golang
-   gRPC
-   gorm
-   mysql
-   redis

## Getting Started

### Prerequisites

-   golang 1.21 or [later](https://go.dev)
-   docker
-   makefile

### Installation

1. Clone this repo
2. Copy `config.example.yaml` in `config` and paste it in the same directory with `.example` removed from its name.
3. Run `go mod download` to download all the dependencies.

### Running
1. Run `docker-compose up -d`
2. Run `make server` or `go run ./src/.`

### Testing
1. Run `make test` or `go test  -v -coverpkg ./... -coverprofile coverage.out -covermode count ./...`

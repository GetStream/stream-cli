NAME = stream-cli

$(NAME):
	@go install ./cmd/$(NAME)

build:
	@go build ./cmd/$(NAME)

test:
	@go test -v ./...

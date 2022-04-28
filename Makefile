NAME = stream-cli

GOLANGCI_VERSION = 1.45.0
GOLANGCI = .bin/golangci/$(GOLANGCI_VERSION)/golangci-lint
$(GOLANGCI):
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(dir $(GOLANGCI)) v$(GOLANGCI_VERSION)

lint: $(GOLANGCI) $(NAME)
	$(GOLANGCI) -v run ./...

lint-fix: $(GOLANGCI) $(NAME)
	$(GOLANGCI) -v run --fix ./...

$(NAME):
	@go install ./cmd/$(NAME)

build:
	@go build ./cmd/$(NAME)

test:
	@go test -v ./...

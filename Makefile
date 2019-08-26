.PHONY: lint

lint:
	if [ ! -f ./bin/golangci-lint ] ; \
	then \
		wget -O - -q https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s latest ; \
	fi;
	./bin/golangci-lint run

test: lint
	go test -covermode=count -coverprofile=count.out -v ./...

mock:
	@rm -rf mocks
	mockery -dir . -name RedirectRepository
	mockery -dir . -name RedirectService

build:
	cd cmd/shortener && GOOS=linux go build -ldflags="-s -w" -o ../../shortener

run:
	go run cmd/shortener/main.go

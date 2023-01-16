
all: sync-dependencies test-run fmt lint

# Sync dependencies command
sync-dependencies:
	@echo "Syncing dependencies with go mod tidy:"
	go mod tidy

# Update dependencies command
update-dependencies:
	@echo "Updating dependencies..."
	go get -u ./...
	@echo "Syncing dependencies with go mod tidy..."
	go mod tidy

# Test commands
test-run:
	@go test -v -short -race ./...

test-cover-run:
	@echo "Running tests and generating report:"
	go test ./... -covermode=atomic -coverprofile=/tmp/coverage.out -coverpkg=./... -count=1
	go tool cover -func /tmp/coverage.out | tail -n 1 | awk '{ print "Total coverage: " $$3 }'
	go tool cover -html=/tmp/coverage.out

lint:
	@echo "Executing golangci-lint$(if $(FLAGS), with flags: $(FLAGS))"
	golangci-lint run

fury-code-quality:
	@echo "Executing fury code-quality (required VPN), this execution may take a while because you need to download the docker image of the service, if it fails, make sure your fury client is up to date."
	fury code-quality run

# format command
fmt:
	@echo "Executing go fmt"
	go fmt $(PACKAGES)
	gofumpt -w .

.PHONY: build-test-check fmt fmt-diff lint lint-fix fix fix-diff test update-direct-dependencies update-all-dependencies


build-test-check:
	go test -run='^$$' ./...

test:
	go test -v -count=1 ./...

fmt:
	go fmt ./...

fmt-diff:
	gofmt -d .

lint:
	go tool golangci-lint version
	go tool golangci-lint run

lint-fix:
	go tool golangci-lint version
	go tool golangci-lint run --fix

## go fix is the Go modernizer tool
fix-diff:
	go version
	go fix -diff ./...

fix:
	go version
	go fix ./...

update-direct-dependencies:
	go list -m -f '{{if not .Indirect}}{{.Path}}{{end}}' all | xargs go get -u
	go mod tidy

update-all-dependencies:
	go get -u ./...
	go mod tidy
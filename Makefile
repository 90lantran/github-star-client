build:
	@echo "build client ..."
	go build -o client cmd/*.go

unit-test:
	cd pkg/client && go test ./... -coverprofile cover.out

code-coverage: unit-test
	cd pkg/client && go tool cover -html=cover.out




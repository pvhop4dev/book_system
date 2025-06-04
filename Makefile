run-api:
	go run cmd/main.go

run-test:
	go run cmd/test.go

gen-doc:
	swag init --parseDependency --parseInternal -g cmd/main.go

check-vuln:
	go install golang.org/x/vuln/cmd/govulncheck@latest
	govulncheck ./...
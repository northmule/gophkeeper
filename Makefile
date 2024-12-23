GOBIN ?= $$(go env GOPATH)/bin

.PHONY: install-go-test-coverage
install-go-test-coverage:
	go install github.com/vladopajic/go-test-coverage/v2@latest

.PHONY: check-coverage
check-coverage: install-go-test-coverage
	go test ./... -coverprofile=./cover.out -covermode=atomic -coverpkg=./...
	${GOBIN}/go-test-coverage --config=./.testcoverage.yml --badge-file-name=./badges/cover.svg


.PHONY: start-db
start-db:
	cd ./build/package/docker/postgres && docker compose up -d

.PHONY: cover-html
run-analyzer:
	go tool cover -html=cover.out
.PHONY: help unit_test integration_test test lint coverage_report cpu_profile mem_profile

help:
	cat Makefile

unit_test:
	go test -v ./internal/...

integration_test:
	go test -v ./integration_test/...

test: unit_test integration_test

lint:
	go fmt ./...
	find . -name '*.go' -exec goimports -local application-design/ -w {} +
	find . -name '*.go' -exec golines -w {} -m 120 \;
	swag fmt ./...
	golangci-lint run ./...


coverage_report:
	go test -coverpkg=./... -count=1 -coverprofile=.coverage.out ./...
	go tool cover -html .coverage.out -o .coverage.html
	open ./.coverage.html

cpu_profile:
	go test -cpuprofile=profiles/cpu.prof  ./integration_test
	go tool pprof -http=:6061 profiles/cpu.prof

mem_profile:
	go test -memprofile=profiles/mem.prof ./integration_test
	go tool pprof -http=:6061 profiles/mem.prof

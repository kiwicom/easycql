all: test

clean:
	rm -rf bin
	rm -rf tests/*_easycql.go
	rm -rf benchmark/*_easycql.go

build:
	go build -i -o ./bin/easycql ./easycql

generate: build
	bin/easycql -stubs \
		./tests/cql_types.go \
		./tests/data.go \
		./tests/nothing.go \

	bin/easycql -all ./tests/cql_types.go
	bin/easycql -all ./tests/data.go
	bin/easycql -all ./tests/nothing.go

test: generate
	go test \
		./tests \
		./gen

bench: build
	./bench.sh

lint:
	golangci-lint run --max-same-issues 0 --max-issues-per-linter 0

.PHONY: clean generate test build

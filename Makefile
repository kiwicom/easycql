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
	#cd benchmark && go test -benchmem -tags use_easycql -bench .
	#golint -set_exit_status ./tests/*_easycql.go

bench: build
	./bench.sh

.PHONY: clean generate test build

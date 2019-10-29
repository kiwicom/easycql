#!/bin/sh
# Run benchmarks with different setups

# Plain gocql
rm -f ./tests/benchmark_easycql.go
echo "---- plain gocql ----"
go test "-bench=.*" ./tests | tee bench-gocql.txt

# easycql conservative mode
echo
echo "---- easycql conservative ----"
bin/easycql -conservative -all ./tests/benchmark.go
go test "-bench=.*" ./tests | tee bench-conservative.txt

# easycql optimized mode
echo
echo "---- easycql optimized ----"
bin/easycql -all ./tests/benchmark.go
go test "-bench=.*" ./tests | tee bench-optimized.txt

# Compare results
echo
echo "---- comparison ----"
echo

if ! command -v benchcmp >/dev/null
then
  echo "benchcmp not installed, please run"
  echo "  GO111MODULE=off go get -u golang.org/x/tools/cmd/benchcmp"
  echo "to display benchmark comparison here"
  exit 1
fi

echo plain gocql vs easycql conservative:
benchcmp bench-gocql.txt bench-conservative.txt
echo
echo plain gocql vs easycql optimized:
benchcmp bench-gocql.txt bench-optimized.txt
echo
echo easycql conservative vs easycql optimized:
benchcmp bench-conservative.txt bench-optimized.txt

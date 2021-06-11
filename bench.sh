#!/bin/sh
# Run benchmarks with different setups

rm bench-gocql.txt bench-conservative.txt bench-optimized.txt

for i in {1..5};
do

	# Plain gocql
	rm -f ./tests/benchmark_easycql.go
	echo "---- plain gocql ----"
	go test "-bench=.*" ./tests -benchmem | tee -a bench-gocql.txt

	# easycql conservative mode
	echo
	echo "---- easycql conservative ----"
	bin/easycql -conservative -all ./tests/benchmark.go
	go test "-bench=.*" ./tests -benchmem | tee -a bench-conservative.txt

	# easycql optimized mode
	echo
	echo "---- easycql optimized ----"
	bin/easycql -all ./tests/benchmark.go
	go test "-bench=.*" ./tests -benchmem | tee -a bench-optimized.txt
done

# Compare results
echo
echo "---- comparison ----"
echo

if ! command -v benchstat >/dev/null
then
  echo "benchstat not installed, please run"
  echo "  go install golang.org/x/perf/cmd/benchstat"
  echo "to display benchmark comparison here"
  exit 1
fi

echo plain gocql vs easycql conservative:
benchstat bench-gocql.txt bench-conservative.txt
echo
echo plain gocql vs easycql optimized:
benchstat bench-gocql.txt bench-optimized.txt
echo
echo easycql conservative vs easycql optimized:
benchstat bench-conservative.txt bench-optimized.txt

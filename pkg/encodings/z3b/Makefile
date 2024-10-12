.PHONY: all test bench bench-profiler clean

all: test bench

test:
	grc go test -tags prod -v ./...

bench: cpu.prof
	grc go test -run=^$ -benchmem -bench=.  ./...

bench-profiler:
	rm -f cpu.prof mem.prof
	grc go test -run=^$ -benchmem -cpuprofile=cpu.prof -memprofile=mem.prof -bench=.  ./...

clean:
	rm -f cpu.prof mem.prof

bench-z3b:
	rm -f cpu.prof mem.prof
	grc go test -run=^$ -benchmem -cpuprofile=cpu.prof -memprofile=mem.prof -bench='.*Z3b.*'  ./...

mem.prof: bench-z3b

cpu.prof: bench-z3b

profile-z3b-mem: mem.prof
	go tool pprof mem.prof

profile-z3b-cpu: cpu.prof
	go tool pprof cpu.prof

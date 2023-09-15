[![Go Reference](https://pkg.go.dev/badge/github.com/rekby/safemutex.svg)](https://pkg.go.dev/github.com/rekby/safemutex)
[![Coverage Status](https://coveralls.io/repos/github/rekby/safemutex/badge.svg?branch=master)](https://coveralls.io/github/rekby/safemutex?branch=master)
[![GoReportCard](https://goreportcard.com/badge/github.com/rekby/safemutex)](https://goreportcard.com/report/github.com/rekby/safemutex)

# Safe mutex

The package inspired by [Rust mutex](https://doc.rust-lang.org/std/sync/struct.Mutex.html). 

Main idea: mutex contains guarded data and no way to use the data with unlocked mutex.

get command:
```bash
go get github.com/rekby/safemutex
```

Example:
```go
package main

import (
	"fmt"
	"github.com/rekby/safemutex"
)

type GuardedStruct struct {
	Name string
	Val  int
}

func main() {
	simleIntMutex := safemutex.New(1)
	simleIntMutex.Lock(func(synced int) int {
		fmt.Println(synced)
		return synced
	})

	mutexWithStruct := safemutex.New(GuardedStruct{Name: "test", Val: 1})
	mutexWithStruct.Lock(func(synced GuardedStruct) GuardedStruct {
		fmt.Println(synced)
		return synced
	})
}
```


# Benchmark result

Safe mutex performance depends on size of stored structure. Benchmarks was with one int value.
If safe mutex in hot way and has a large structure - use MutexWithPointers or RWMutexWithPointers and store
a pointer to structure in mutex. It will some reduce guarantees, but will better performance.  

```
Macbook M1

BenchmarkSyncMutexLock-10                  	88132146	        13.57 ns/op	       0 B/op	       0 allocs/op
BenchmarkSafeMutexLock-10                  	88716652	        13.56 ns/op	       0 B/op	       0 allocs/op
BenchmarkSafeMutexWithPointersLock-10      	87819339	        13.64 ns/op	       0 B/op	       0 allocs/op
BenchmarkSyncRWMutexLock-10                	64879916	        18.52 ns/op	       0 B/op	       0 allocs/op
BenchmarkSafeRWMutexLock-10                	64612960	        18.50 ns/op	       0 B/op	       0 allocs/op
BenchmarkSafeRWMutexWithPointersLock-10    	64686685	        18.58 ns/op	       0 B/op	       0 allocs/op

Macbook Intel
BenchmarkSyncMutexLock-12                       100000000               10.94 ns/op            0 B/op          0 allocs/op
BenchmarkSafeMutexLock-12                       64584142                18.90 ns/op            0 B/op          0 allocs/op
BenchmarkSafeMutexWithPointersLock-12           76630285                15.79 ns/op            0 B/op          0 allocs/op
BenchmarkSyncRWMutexLock-12                     50869114                22.39 ns/op            0 B/op          0 allocs/op
BenchmarkSafeRWMutexLock-12                     41917308                28.80 ns/op            0 B/op          0 allocs/op
BenchmarkSafeRWMutexWithPointersLock-12         44739602                26.96 ns/op            0 B/op          0 allocs/op

```

[![Go Reference](https://pkg.go.dev/badge/github.com/rekby/safemutex.svg)](https://pkg.go.dev/github.com/rekby/safemutex)
[![Coverage Status](https://coveralls.io/repos/github/rekby/safe-mutex/badge.svg?branch=master)](https://coveralls.io/github/rekby/safe-mutex?branch=master)
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

```
BenchmarkSyncMutexLock
BenchmarkSyncMutexLock-10                	85375321	        13.61 ns/op	       0 B/op	       0 allocs/op
BenchmarkSyncRWMutexLock
BenchmarkSyncRWMutexLock-10              	64748929	        18.56 ns/op	       0 B/op	       0 allocs/op
BenchmarkSafeMutexLock
BenchmarkSafeMutexLock-10                	86766590	        13.79 ns/op	       0 B/op	       0 allocs/op
BenchmarkSafeMutexWithPointersLock
BenchmarkSafeMutexWithPointersLock-10    	87910747	        13.66 ns/op	       0 B/op	       0 allocs/op
```

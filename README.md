[![Go Reference](https://pkg.go.dev/badge/github.com/rekby/go-safe-mutex.svg)](https://pkg.go.dev/github.com/rekby/go-safe-mutex)
[![Coverage Status](https://coveralls.io/repos/github/rekby/go-safe-mutex/badge.svg?branch=master)](https://coveralls.io/github/rekby/go-safe-mutex?branch=master)
[![GoReportCard](https://goreportcard.com/badge/github.com/rekby/go-safe-mutex)](https://goreportcard.com/report/github.com/rekby/go-safe-mutex)

# Safe mutex

The package inspired by [Rust mutex](https://doc.rust-lang.org/std/sync/struct.Mutex.html). 

Main idea: mutex contains guarded data and no way to use the data with unlocked mutex.


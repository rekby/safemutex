[![Go Reference](https://pkg.go.dev/badge/github.com/rekby/safe-mutex.svg)](https://pkg.go.dev/github.com/rekby/safe-mutex)
[![Coverage Status](https://coveralls.io/repos/github/rekby/safe-mutex/badge.svg?branch=master)](https://coveralls.io/github/rekby/safe-mutex?branch=master)
[![GoReportCard](https://goreportcard.com/badge/github.com/rekby/safe-mutex)](https://goreportcard.com/report/github.com/rekby/safe-mutex)

# Safe mutex

The package inspired by [Rust mutex](https://doc.rust-lang.org/std/sync/struct.Mutex.html). 

Main idea: mutex contains guarded data and no way to use the data with unlocked mutex.


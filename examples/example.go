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
	simleIntMutex.Lock(func(value int) (newValue int) {
		fmt.Println(value)
		return value
	})

	mutexWithStruct := safemutex.New(GuardedStruct{Name: "test", Val: 1})
	mutexWithStruct.Lock(func(value GuardedStruct) (newValue GuardedStruct) {
		fmt.Println(value)
		return value
	})
}

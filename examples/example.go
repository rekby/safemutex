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

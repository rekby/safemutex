package safemutex_test

import (
	"fmt"
	"github.com/rekby/safemutex"
)

func ExampleNew() {
	type Struct struct {
		m safemutex.Mutex[int]
	}

	s := Struct{m: safemutex.New(1)}

	s.m.Lock(func(synced int) int {
		return synced + 1
	})

	var mutexVal int
	s.m.Lock(func(synced int) int {
		mutexVal = synced
		return synced
	})
	fmt.Println("mutexVal:", mutexVal)

	// Output:
	// mutexVal: 2
}

func ExampleNew_with_struct() {
	type GuargedStruct struct {
		Name string
		Val  int
	}
	type Struct struct {
		m safemutex.Mutex[GuargedStruct]
	}

	s := Struct{m: safemutex.New(GuargedStruct{Name: "test-name", Val: 15})}

	s.m.Lock(func(synced GuargedStruct) GuargedStruct {
		synced.Val += 1
		return synced
	})

	var mutexVal GuargedStruct
	s.m.Lock(func(synced GuargedStruct) GuargedStruct {
		mutexVal = synced
		return synced
	})
	fmt.Printf("name: %q, val: %v", mutexVal.Name, mutexVal.Val)

	// Output:
	// name: "test-name", val: 16
}

func ExampleMutex_Lock() {
	type Struct struct {
		m safemutex.Mutex[int]
	}

	var s Struct

	s.m.Lock(func(synced int) int {
		return synced + 1
	})

	var mutexVal int
	s.m.Lock(func(synced int) int {
		mutexVal = synced
		return synced
	})
	fmt.Println("mutexVal:", mutexVal)

	// Output:
	// mutexVal: 1
}

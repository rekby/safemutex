package safemutex_test

import (
	"fmt"
	"github.com/rekby/safemutex"
)

func ExampleNewWithPointers() {
	type Struct struct {
		m safemutex.MutexWithPointers[map[string]int]
	}

	s := Struct{m: safemutex.NewWithPointers(map[string]int{})}

	s.m.Lock(func(synced map[string]int) map[string]int {
		synced["asd"] = 123
		return synced
	})

	var mutexVal map[string]int
	s.m.Lock(func(synced map[string]int) map[string]int {
		mutexVal = synced
		return synced
	})
	fmt.Println("mutexVal:", mutexVal)

	// Output:
	// mutexVal: map[asd:123]
}

func ExampleNewWithPointers_with_struct() {
	type GuargedStruct struct {
		Name string
		Val  *int
	}
	type Struct struct {
		m safemutex.MutexWithPointers[GuargedStruct]
	}

	s := Struct{m: safemutex.NewWithPointers(GuargedStruct{Name: "test-name"})}

	s.m.Lock(func(synced GuargedStruct) GuargedStruct {
		if synced.Val == nil {
			v := 15
			synced.Val = &v
		}
		return synced
	})

	var mutexVal GuargedStruct
	s.m.Lock(func(synced GuargedStruct) GuargedStruct {
		*synced.Val += 1
		mutexVal = synced
		return synced
	})
	fmt.Printf("name: %q, val: %v", mutexVal.Name, *mutexVal.Val)

	// Output:
	// name: "test-name", val: 16
}

func ExampleMutexWithPointers_Lock() {
	var m safemutex.MutexWithPointers[map[string]int]

	m.Lock(func(synced map[string]int) map[string]int {
		if synced == nil {
			synced = map[string]int{}
		}
		synced["asd"] = 22
		return synced
	})

	var mutexVal map[string]int
	m.Lock(func(synced map[string]int) map[string]int {
		mutexVal = synced
		return synced
	})
	fmt.Println("mutexVal:", mutexVal)

	// Output:
	// mutexVal: map[asd:22]
}

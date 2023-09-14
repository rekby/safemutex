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

	s.m.Lock(func(value int) (newValue int) {
		return value + 1
	})

	var mutexVal int
	s.m.Lock(func(value int) (newValue int) {
		mutexVal = value
		return value
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

	s.m.Lock(func(value GuargedStruct) (newValue GuargedStruct) {
		value.Val += 1
		return value
	})

	var mutexVal GuargedStruct
	s.m.Lock(func(value GuargedStruct) (newValue GuargedStruct) {
		mutexVal = value
		return value
	})
	fmt.Printf("name: %q, val: %v", mutexVal.Name, mutexVal.Val)

	// Output:
	// name: "test-name", val: 16
}

func ExampleNewWithPointers() {
	val1 := 1
	val2 := 2

	m := safemutex.NewWithPointers(&val1)

	m.Lock(func(value *int) (newValue *int) {
		fmt.Println(*value)
		return &val2
	})

	m.Lock(func(value *int) (newValue *int) {
		fmt.Println(*value)
		return value
	})

	// Output:
	// 1
	// 2
}

func ExampleMutex_Lock() {
	type Struct struct {
		m safemutex.Mutex[int]
	}

	var s Struct

	s.m.Lock(func(value int) (newValue int) {
		return value + 1
	})

	var mutexVal int
	s.m.Lock(func(value int) (newValue int) {
		mutexVal = value
		return value
	})
	fmt.Println("mutexVal:", mutexVal)

	// Output:
	// mutexVal: 1
}

//go:build go1.19
// +build go1.19

package safemutex_test

import (
	"fmt"
	"github.com/rekby/safemutex"
)

func ExampleMutex_TryLock() {
	m := safemutex.New(1)

	var (
		outerLock bool
		innerLock bool
	)

	outerLock = m.TryLock(func(synced int) int {
		innerLock = m.TryLock(func(innerSynced int) int {
			innerLock = true
			innerSynced += 1
			fmt.Println(innerSynced)
			return innerSynced
		})
		fmt.Println(synced)
		return synced
	})
	fmt.Println("inner lock:", innerLock)
	fmt.Println("outer lock:", outerLock)

	// Output:
	// 1
	// inner lock: false
	// outer lock: true
}

func ExampleMutexWithPointers_TryLock() {
	m := safemutex.NewWithPointers(map[string]int{})

	var (
		outerLock bool
		innerLock bool
	)

	outerLock = m.TryLock(func(synced map[string]int) map[string]int {
		innerLock = m.TryLock(func(innerSynced map[string]int) map[string]int {
			innerLock = true
			innerSynced["inner"] = 222
			fmt.Println(innerSynced)
			return innerSynced
		})
		synced["outer"] = 123
		fmt.Println(synced)
		return synced
	})
	fmt.Println("inner lock:", innerLock)
	fmt.Println("outer lock:", outerLock)

	// Output:
	// map[outer:123]
	// inner lock: false
	// outer lock: true

}

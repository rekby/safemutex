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

	outerLock = m.TryLock(func(value int) (newValue int) {
		innerLock = m.TryLock(func(value int) (newValue int) {
			innerLock = true
			value += 1
			fmt.Println(value)
			return value
		})
		fmt.Println(value)
		return value
	})
	fmt.Println("inner lock:", innerLock)
	fmt.Println("outer lock:", outerLock)

	// Output:
	// 1
	// inner lock: false
	// outer lock: true
}

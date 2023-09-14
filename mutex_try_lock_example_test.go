package safe_mutex_test

import (
	"fmt"
	safe_mutex "safe-mutex"
)

func ExampleMutex_TryLock() {
	m := safe_mutex.New(1)

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

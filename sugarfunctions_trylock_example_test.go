//go:build go1.19
// +build go1.19

package safemutex_test

import (
	"fmt"

	"github.com/rekby/safemutex"
)

func ExampleTryLockWithResult() {
	counter := Counter{m: safemutex.RWNew(1)}

	mess, ok := safemutex.TryLockWithResult(&counter.m, func(state int) (newState int, result string) {
		return state + 1, fmt.Sprint("Last value of counter: ", state)
	})
	fmt.Println("Ok:", ok)
	fmt.Println(mess)

	locked := make(chan bool)
	allowFree := make(chan bool)
	unlocked := make(chan bool)
	go func() {
		counter.m.Lock(func(synced int) int {
			close(locked)
			<-allowFree
			return synced
		})
		close(unlocked)
	}()

	<-locked
	_, ok = safemutex.TryLockWithResult(&counter.m, func(synced int) (int, int) {
		panic("the callback will not be called")
	})

	fmt.Println("False when mutex already locked:", ok)

	close(allowFree)
	<-unlocked

	// Output:
	// Ok: true
	// Last value of counter: 1
	// False when mutex already locked: false
}

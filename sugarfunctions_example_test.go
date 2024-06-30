package safemutex_test

import (
	"fmt"

	"github.com/rekby/safemutex"
)

type Counter struct {
	m safemutex.RWMutex[int]
}

func ExampleLockWithResult() {
	counter := Counter{m: safemutex.RWNew(1)}

	mess := safemutex.LockWithResult(&counter.m, func(state int) (newState int, result string) {
		return state + 1, fmt.Sprint("Last value of counter: ", state)
	})
	fmt.Println(mess)

	counter.m.RLock(func(synced int) {
		fmt.Println("Internal state:", synced)
	})

	// Output:
	// Last value of counter: 1
	// Internal state: 2
}

func ExampleRLockWithResult() {
	counter := Counter{m: safemutex.RWNew(5)}

	mess := safemutex.RLockWithResult(&counter.m, func(synced int) (result string) {
		return fmt.Sprint("Double value: ", synced*2)
	})
	fmt.Println(mess)

	counter.m.RLock(func(synced int) {
		fmt.Println("Internal state:", synced)
	})

	// Output:
	// Double value: 10
	// Internal state: 5
}

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

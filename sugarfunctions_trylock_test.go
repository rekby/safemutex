//go:build go1.19
// +build go1.19

package safemutex

import (
	"fmt"
	"testing"
)

func TestTryLock(t *testing.T) {
	t.Run("check_try", func(t *testing.T) {
		m := New(1)
		res, ok := TryLockWithResult(&m, func(synced int) (int, string) {
			return synced + 1, fmt.Sprintf("old: %v", synced)
		})
		if !ok {
			t.Fatal()
		}
		if res != "old: 1" {
			t.Fatal(res)
		}

		locked := make(chan bool)
		allowFree := make(chan bool)
		freed := make(chan bool)
		var backgroundLockResult bool
		go func() {
			res, backgroundLockResult = TryLockWithResult(&m, func(state int) (int, string) {
				close(locked)
				<-allowFree
				return state + 2, fmt.Sprintf("old2: %v", state)
			})
			close(freed)
		}()

		<-locked

		_, ok = TryLockWithResult(&m, func(state int) (int, string) {
			t.Fatal("the callback must not be called")
			return -1, ""
		})

		if ok {
			t.Fatal()
		}

		close(allowFree)
		<-freed
		if !backgroundLockResult {
			t.Fatal()
		}
	})

}

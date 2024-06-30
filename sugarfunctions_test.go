package safemutex

import (
	"fmt"
	"testing"
)

func TestLockWithResult(t *testing.T) {
	t.Run("mutex", func(t *testing.T) {
		m := New(1)
		res := LockWithResult(&m, func(synced int) (int, string) {
			return synced + 2, fmt.Sprintf("old: %v", synced)
		})

		if res != "old: 1" {
			t.Fatal(res)
		}
	})
	t.Run("rw mutex", func(t *testing.T) {
		m := RWNew(1)
		res := LockWithResult(&m, func(synced int) (int, string) {
			return synced + 2, fmt.Sprintf("old: %v", synced)
		})

		if res != "old: 1" {
			t.Fatal(res)
		}
	})
	t.Run("rw mutex with pointers", func(t *testing.T) {
		m := RWNewWithPointers(1)
		res := LockWithResult(&m, func(synced int) (int, string) {
			return synced + 2, fmt.Sprintf("old: %v", synced)
		})

		if res != "old: 1" {
			t.Fatal(res)
		}
	})
}

func TestRLockWithResult(t *testing.T) {
	t.Run("rw_mutex", func(t *testing.T) {
		m := RWNew(1)
		res := RLockWithResult(&m, func(synced int) string {
			return fmt.Sprintf("old: %v", synced)
		})

		if res != "old: 1" {
			t.Fatal(res)
		}
	})
	t.Run("rw_mutex_with_pointers", func(t *testing.T) {
		m := RWNewWithPointers(1)
		res := RLockWithResult(&m, func(synced int) string {
			return fmt.Sprintf("old: %v", synced)
		})

		if res != "old: 1" {
			t.Fatal(res)
		}
	})
}

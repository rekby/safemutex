//go:build go1.19
// +build go1.19

package safemutex

import "testing"

func TestMutexTryLock(t *testing.T) {
	callCount := 0
	initialValue := -1
	tmpValue := -2
	targetValue := 123
	m := New(initialValue)
	lockedOuter := m.TryLock(func(synced int) int {
		callCount++

		lockedInner := m.TryLock(func(synced int) int {
			t.Fatal()
			return tmpValue
		})
		if lockedInner {
			t.Fatal()
		}
		return targetValue
	})
	if !lockedOuter {
		t.Fatal()
	}
	if m.value != targetValue {
		t.Fatal(m.value)
	}
}

func TestMutexWithPointersTryLock(t *testing.T) {
	callCount := 0
	key := "test"
	initialValue := -1
	tmpValue := -2
	targetValue := 123
	m := NewWithPointers(map[string]int{key: initialValue})
	lockedOuter := m.TryLock(func(synced map[string]int) map[string]int {
		callCount++

		lockedInner := m.TryLock(func(synced map[string]int) map[string]int {
			t.Fatal()
			synced[key] = tmpValue
			return synced
		})
		if lockedInner {
			t.Fatal()
		}
		synced[key] = targetValue
		return synced
	})
	if !lockedOuter {
		t.Fatal()
	}
	if m.value[key] != targetValue {
		t.Fatal(m.value)
	}
}

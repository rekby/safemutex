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

func TestRWMutexTryLock(t *testing.T) {
	callCount := 0
	initialValue := -1
	tmpValue := -2
	targetValue := 123
	m := RWNew(initialValue)
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

func TestRWMutexTryRLock(t *testing.T) {
	callCount := 0
	initialValue := -1
	targetValue := 123
	m := RWNew(initialValue)

	// try rlock into lock
	lockedOuter := m.TryLock(func(synced int) int {
		callCount++

		lockedInner := m.TryRLock(func(synced int) {
			callCount++
			t.Fatal()
		})
		if lockedInner {
			t.Fatal()
		}
		return targetValue
	})
	if !lockedOuter {
		t.Fatal()
	}

	if callCount != 1 {
		t.Fatal(callCount)
	}

	if m.value != targetValue {
		t.Fatal(m.value)
	}

	// try rlock into rlock
	callCount = 0
	lockedOuter = m.TryRLock(func(synced int) {
		callCount++
		innerLock := m.TryRLock(func(synced int) {
			callCount++
		})
		if !innerLock {
			t.Fatal()
		}
	})
	if !lockedOuter {
		t.Fatal()
	}
	if callCount != 2 {
		t.Fatal(callCount)
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

func TestRWMutexWithPointersTryLock(t *testing.T) {
	callCount := 0
	key := "test"
	initialValue := -1
	tmpValue := -2
	targetValue := 123
	m := RWNewWithPointers(map[string]int{key: initialValue})
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

func TestRWMutexWithPointersTryRLock(t *testing.T) {
	callCount := 0
	key := "test"
	initialValue := -1
	targetValue := 123
	m := RWNewWithPointers(map[string]int{key: initialValue})

	// try rlock into lock
	lockedOuter := m.TryLock(func(synced map[string]int) map[string]int {
		callCount++

		lockedInner := m.TryRLock(func(synced map[string]int) {
			callCount++
			t.Fatal()
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

	if callCount != 1 {
		t.Fatal(callCount)
	}

	if m.value[key] != targetValue {
		t.Fatal(m.value)
	}

	// try rlock into rlock
	callCount = 0
	lockedOuter = m.TryRLock(func(synced map[string]int) {
		callCount++
		innerLock := m.TryRLock(func(synced map[string]int) {
			callCount++
		})
		if !innerLock {
			t.Fatal()
		}
	})
	if !lockedOuter {
		t.Fatal()
	}
	if callCount != 2 {
		t.Fatal(callCount)
	}
}

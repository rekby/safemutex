//go:build go1.19
// +build go1.19

package safe_mutex

import "testing"

func TestTryLock(t *testing.T) {
	callCount := 0
	initialValue := -1
	tmpValue := -2
	targetValue := 123
	m := New(initialValue)
	lockedOuter := m.TryLock(func(value int) (newValue int) {
		callCount++

		lockedInner := m.TryLock(func(value int) (newValue int) {
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

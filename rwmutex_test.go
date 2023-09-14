package safemutex

import (
	"errors"
	"testing"
	"time"
)

func TestRWMutexZero(t *testing.T) {
	targetValue := 123

	var m RWMutex[int]
	callCount := 0
	m.Lock(func(synced int) int {
		callCount++
		return targetValue
	})

	if callCount != 1 {
		t.Fatal(callCount)
	}
	if m.value != targetValue {
		t.Fatal(m.value)
	}
}

func TestRWMutexLocked(t *testing.T) {
	tmpValue := -1
	targetValue := 123

	var m RWMutex[int]
	callCount := 0
	innerCompleted := make(chan bool)
	m.Lock(func(synced int) int {
		callCount++

		go func() {
			m.Lock(func(innerSynced int) int {
				callCount++
				return targetValue
			})
			close(innerCompleted)
		}()

		time.Sleep(time.Millisecond)
		return tmpValue
	})

	<-innerCompleted
	if callCount != 2 {
		t.Fatal(callCount)
	}
	if m.value != targetValue {
		t.Fatal(m.value)
	}
}

func TestRWMutexReadLocked(t *testing.T) {
	targetValue := 123

	var m RWMutexWithPointers[int]

	called := 0
	m.Lock(func(synced int) int {
		called++
		return targetValue
	})

	m.RLock(func(synced int) {
		called++
		innerCompleted := make(chan bool)
		m.RLock(func(int) {
			close(innerCompleted)
		})
		<-innerCompleted
	})

	if called != 2 {
		t.Fatal(called)
	}

	if m.value != targetValue {
		t.Fatal(m.value)
	}
}

func TestRWNew(t *testing.T) {
	t.Run("WithPointers", func(t *testing.T) {
		defer func() {
			err := recover()
			if !errors.Is(err.(error), errContainPointers) {
				t.Fatal(err)
			}
		}()

		_ = RWNew(struct{ v *int }{})
	})
	t.Run("WithoutPointers", func(t *testing.T) {
		_ = RWNew(struct{ v int }{})
	})
}

func TestRWMutexInitialized(t *testing.T) {
	t.Run("New", func(t *testing.T) {
		m := RWNew(123)
		if !m.initialized {
			t.Fatal()
		}
	})
	t.Run("Lock", func(t *testing.T) {
		var m RWMutex[int]
		if m.initialized {
			t.Fatal()
		}
		m.Lock(func(value int) (newValue int) {
			return value
		})
		if !m.initialized {
			t.Fatal()
		}
	})

	t.Run("RLock", func(t *testing.T) {
		var m RWMutex[int]
		if m.initialized {
			t.Fatal()
		}
		m.RLock(func(value int) {
			// stub
		})
		if !m.initialized {
			t.Fatal()
		}
	})

	t.Run("BadValueLock", func(t *testing.T) {
		var m RWMutex[*int]
		if m.initialized {
			t.Fatal()
		}

		defer func() {
			err := recover().(error)
			if !errors.Is(err, errContainPointers) {
				t.Fatal(err)
			}
		}()

		m.Lock(func(value *int) (newValue *int) {
			t.Fatal()
			return nil
		})
	})
}

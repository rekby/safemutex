package safemutex

import (
	"errors"
	"testing"
)

func TestMutexPoisoned(t *testing.T) {
	t.Run("Default", func(t *testing.T) {
		initialValue := 123
		secondValue := -1
		targetValue := initialValue

		m := New(initialValue)

		defer func() {
			err := recover().(error)
			if !errors.Is(err, ErrPoisoned) {
				t.Fatal(err)
			}
			if m.value != targetValue {
				t.Fatal(m.value)
			}
		}()

		hasPanic := false

		// panic in mutex
		func() {
			defer func() {
				if recover() != nil {
					hasPanic = true
				}
			}()

			m.Lock(func(value int) (newValue int) {
				panic("test")
			})
		}()

		if !hasPanic {
			t.Fatal()
		}

		m.Lock(func(value int) (newValue int) {
			return secondValue
		})
	})
}

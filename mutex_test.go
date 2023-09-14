package safemutex

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

func TestZeroOk(t *testing.T) {
	targetValue := 123

	var m Mutex[int]
	callCount := 0
	m.Lock(func(value int) (newValue int) {
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

func TestLockedOk(t *testing.T) {
	tmpValue := -1
	targetValue := 123

	var m Mutex[int]
	callCount := 0
	innerCompleted := make(chan bool)
	m.Lock(func(value int) (newValue int) {
		callCount++

		go func() {
			m.Lock(func(value int) (newValue int) {
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

func TestValueWithPointers(t *testing.T) {
	t.Run("Default", func(t *testing.T) {
		defer func() {
			err := recover()
			if err != errContainPointers {
				t.Fatal(err)
			}
		}()

		_ = New(struct{ v *int }{})
	})
	t.Run("AllowPointers", func(t *testing.T) {
		_ = NewWithOptions(struct{ v *int }{}, MutexOptions{AllowPointers: true})
	})
}

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
	t.Run("WithAllowPoisoned", func(t *testing.T) {
		initialValue := -1
		secondValue := 123
		targetValue := secondValue

		m := NewWithOptions(initialValue, MutexOptions{AllowPoisoned: true})

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

		if m.value != targetValue {
			t.Fatal(m.value)
		}
	})
}

func TestCheckPointers(t *testing.T) {
	testVal := 0
	tests := []struct {
		name           string
		value          any
		expectedResult bool
	}{
		// base types without pointer
		{
			name:           "int",
			value:          int(0),
			expectedResult: false,
		},
		{
			name:           "string",
			value:          "asd",
			expectedResult: false,
		},
		// base types with pointers
		{
			name:           "pointer",
			value:          &testVal,
			expectedResult: true,
		},
		{
			name:           "nil",
			value:          nil,
			expectedResult: true,
		},
		{
			name:           "slice",
			value:          []int{1, 2, 3},
			expectedResult: true,
		},
		{
			name:           "map",
			value:          map[int]int{1: 1},
			expectedResult: true,
		},
		{
			name:           "map",
			value:          map[int]int{1: 1},
			expectedResult: true,
		},
		{
			name:           "map",
			value:          map[int]int{1: 1},
			expectedResult: true,
		},
		{
			name:           "func",
			value:          func() {},
			expectedResult: true,
		},

		// compound types
		{
			name:           "array without pointers",
			value:          [1]int{0},
			expectedResult: false,
		},
		{
			name:           "array with pointers",
			value:          [1]*int{&testVal},
			expectedResult: true,
		},
		{
			name:           "struct without pointer",
			value:          struct{ val int }{1},
			expectedResult: false,
		},
		{
			name:           "struct with pointer",
			value:          struct{ val *int }{&testVal},
			expectedResult: true,
		},
		{
			name:           "array of struct without pointer",
			value:          [1]struct{ val int }{{1}},
			expectedResult: false,
		},
		{
			name:           "array of struct with pointer",
			value:          [1]struct{ val *int }{{&testVal}},
			expectedResult: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res := checkTypeCanContainPointers(reflect.TypeOf(test.value))
			if res != test.expectedResult {
				t.Fatalf("value %q, expected: %v actual: %v", test.value, test.expectedResult, res)
			}
		})
	}
}

func TestMutexInitialized(t *testing.T) {
	t.Run("New", func(t *testing.T) {
		m := New(123)
		if !m.initialized {
			t.Fatal()
		}
	})
	t.Run("Lock", func(t *testing.T) {
		var m Mutex[int]
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

	t.Run("BadValueLock", func(t *testing.T) {
		var m Mutex[*int]
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

func TestCheckPointersWithUnknownKind(t *testing.T) {
	if !checkTypeCanContainPointers(typeWithInvalidKind{}) {
		t.Fatal()
	}
}

type typeWithInvalidKind struct{}

func (t typeWithInvalidKind) Kind() reflect.Kind {
	return reflect.Invalid
}

func (t typeWithInvalidKind) NumField() int {
	panic("no need")
}

func (t typeWithInvalidKind) Field(_ int) reflect.StructField {
	panic("no need")
}

func (t typeWithInvalidKind) Elem() reflect.Type {
	panic("no need")
}

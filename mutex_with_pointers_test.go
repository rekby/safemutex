package safemutex

import (
	"testing"
	"time"
)

func TestMutexWithPointersZero(t *testing.T) {
	targetValue := 123

	var m MutexWithPointers[map[string]int]
	callCount := 0
	m.Lock(func(synced map[string]int) map[string]int {
		callCount++
		if synced == nil {
			synced = map[string]int{}
		}

		synced["test"] = targetValue
		return synced
	})

	if callCount != 1 {
		t.Fatal(callCount)
	}
	if m.value["test"] != targetValue {
		t.Fatal(m.value)
	}
}

func TestMutexWithPointersLocked(t *testing.T) {
	tmpValue := -1
	targetValue := 123
	key := "test"

	var m MutexWithPointers[map[string]int]
	callCount := 0
	innerCompleted := make(chan bool)
	m.Lock(func(synced map[string]int) map[string]int {
		callCount++

		synced = map[string]int{
			key: tmpValue,
		}

		go func() {
			m.Lock(func(innerSynced map[string]int) map[string]int {
				callCount++
				synced[key] = targetValue
				return synced
			})
			close(innerCompleted)
		}()

		time.Sleep(time.Millisecond)
		return synced
	})

	<-innerCompleted
	if callCount != 2 {
		t.Fatal(callCount)
	}
	if m.value[key] != targetValue {
		t.Fatal(m.value)
	}
}

func TestNewWithPointers(t *testing.T) {
	val := 1
	t.Run("WithPointers", func(t *testing.T) {
		v := NewWithPointers(struct{ v *int }{&val})
		if v.value.v != &val {
			t.Fatal()
		}
	})
	t.Run("WithoutPointers", func(t *testing.T) {
		v := NewWithPointers(struct{ v int }{val})
		if v.value.v != val {
			t.Fatal()
		}
	})
}

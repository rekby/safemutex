package safemutex

import (
	"testing"
	"time"
)

func TestRWMutexWithPointersZeroOk(t *testing.T) {
	targetValue := 123

	var m RWMutexWithPointers[map[string]int]
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

func TestRWMutexWithPointersLocked(t *testing.T) {
	tmpValue := -1
	targetValue := 123
	key := "test"

	var m RWMutexWithPointers[map[string]int]
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

func TestRWMutexWithPointersReadLocked(t *testing.T) {
	targetValue := 123
	key := "test"

	var m RWMutexWithPointers[map[string]int]

	called := 0
	m.Lock(func(synced map[string]int) map[string]int {
		called++
		if synced == nil {
			synced = map[string]int{}
		}
		synced[key] = targetValue
		return synced
	})

	m.RLock(func(synced map[string]int) {
		called++
		innerCompleted := make(chan bool)
		m.RLock(func(synced map[string]int) {
			close(innerCompleted)
		})
		<-innerCompleted
	})

	if called != 2 {
		t.Fatal(called)
	}

	if m.value[key] != targetValue {
		t.Fatal(m.value)
	}
}

func TestRWNewWithPointers(t *testing.T) {
	val := 1
	t.Run("WithPointers", func(t *testing.T) {
		v := RWNewWithPointers(struct{ v *int }{&val})
		if v.value.v != &val {
			t.Fatal()
		}
	})
	t.Run("WithoutPointers", func(t *testing.T) {
		v := RWNewWithPointers(struct{ v int }{val})
		if v.value.v != val {
			t.Fatal()
		}
	})
}

package safemutex

import (
	"fmt"
	"io"
	"sync"
	"testing"
)

func BenchmarkSyncMutexLock(b *testing.B) {
	b.ReportAllocs()

	var val int
	var m sync.Mutex

	m.Lock()
	m.Unlock()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Lock()
		val++
		m.Unlock()
	}
	b.StopTimer()
	fmt.Fprint(io.Discard, val)
}

func BenchmarkSafeMutexLock(b *testing.B) {
	b.ReportAllocs()

	m := New(0)
	m.Lock(func(synced int) int {
		return synced
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Lock(func(synced int) int {
			synced++
			return synced
		})
	}
	b.StopTimer()

	var val int
	m.Lock(func(synced int) int {
		val = synced
		return synced
	})
	fmt.Fprint(io.Discard, val)

}

func BenchmarkSafeMutexWithPointersLock(b *testing.B) {
	b.ReportAllocs()

	m := NewWithPointers(0)
	m.Lock(func(synced int) int {
		return synced
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Lock(func(synced int) int {
			synced++
			return synced
		})
	}
	b.StopTimer()

	var val int
	m.Lock(func(synced int) int {
		val = synced
		return synced
	})
	fmt.Fprint(io.Discard, val)

}

func BenchmarkSyncRWMutexLock(b *testing.B) {
	b.ReportAllocs()

	var val int
	var m sync.RWMutex

	m.Lock()
	m.Unlock()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Lock()
		val++
		m.Unlock()
	}
	b.StopTimer()
	fmt.Fprint(io.Discard, val)
}

func BenchmarkSafeRWMutexLock(b *testing.B) {
	b.ReportAllocs()

	m := RWNew(0)
	m.Lock(func(synced int) int {
		return synced
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Lock(func(synced int) int {
			synced++
			return synced
		})
	}
	b.StopTimer()

	var val int
	m.Lock(func(synced int) int {
		val = synced
		return synced
	})
	fmt.Fprint(io.Discard, val)

}

func BenchmarkSafeRWMutexWithPointersLock(b *testing.B) {
	b.ReportAllocs()

	m := RWNewWithPointers(0)
	m.Lock(func(synced int) int {
		return synced
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Lock(func(synced int) int {
			synced++
			return synced
		})
	}
	b.StopTimer()

	var val int
	m.Lock(func(synced int) int {
		val = synced
		return synced
	})
	fmt.Fprint(io.Discard, val)

}

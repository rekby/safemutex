package safemutex

import "sync"

// RWMutexWithPointers contains guarded value inside, access to value allowed inside callbacks only
// it allow to guarantee not access to the value without lock the mutex
// zero value is usable as mutex with default options and zero value of guarded type
type RWMutexWithPointers[T any] struct {
	mutexBase[T, sync.RWMutex]
}

// RWNewWithPointers create RWMutexWithPointers with initial value and default options.
// RWNewWithPointers call internal checks for T and panic if checks failed, see MutexOptions for details
func RWNewWithPointers[T any](value T) RWMutexWithPointers[T] {
	res := RWMutexWithPointers[T]{
		mutexBase: mutexBase[T, sync.RWMutex]{
			value: value,
		},
	}

	//nolint:govet
	//goland:noinspection GoVetCopyLock
	return res
}

// Lock - call f within locked mutex.
// it will panic with ErrPoisoned if previous locked call exited without return value:
// with panic or runtime.Goexit()
func (m *RWMutexWithPointers[T]) Lock(f ReadWriteCallback[T]) {
	m.m.Lock()
	defer m.m.Unlock()

	m.baseValidateLocked()
	m.callLocked(f)
}

// RLock - call f within locked mutex.
// it will panic with ErrPoisoned if previous locked call exited without return value:
// with panic or runtime.Goexit()
func (m *RWMutexWithPointers[T]) RLock(f ReadCallback[T]) {
	m.m.RLock()
	defer m.m.RUnlock()

	m.baseValidateLocked()
	m.callReadLocked(f)
}

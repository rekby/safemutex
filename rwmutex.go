package safemutex

import (
	"reflect"
	"sync"
)

// RWMutex contains guarded value inside, access to value allowed inside callbacks only
// it allow to guarantee not access to the value without lock the mutex
// zero value is usable as mutex with default options and zero value of guarded type
// RWMutex deny to save value with any type of pointers, which allow accidentally change internal state.
// it will panic if T contains any pointer.
type RWMutex[T any] struct {
	mutexBase[T, sync.RWMutex]
	initOnce    sync.Once
	initialized bool // for tests only
}

// RWNew create RWMutex with initial value and default options.
// RWNew call internal checks for T and panic if checks failed, see MutexOptions for details
func RWNew[T any](value T) RWMutex[T] {
	res := RWMutex[T]{
		mutexBase: mutexBase[T, sync.RWMutex]{
			value: value,
		},
	}

	res.validateLocked()

	//nolint:govet
	//goland:noinspection GoVetCopyLock
	return res
}

// Lock - call f within locked mutex.
// it will panic if value type not pass internal checks
// it will panic with ErrPoisoned if previous locked call exited without return value:
// with panic or runtime.Goexit()
func (m *RWMutex[T]) Lock(f ReadWriteCallback[T]) {
	m.m.Lock()
	defer m.m.Unlock()

	m.validateLocked()
	m.callLocked(f)
}

// RLock - call f within read locked mutex.
// it will panic if value type not pass internal checks
// it will panic with ErrPoisoned if previous locked call exited without return value:
// with panic or runtime.Goexit()
func (m *RWMutex[T]) RLock(f ReadCallback[T]) {
	m.m.Lock()
	defer m.m.Unlock()

	m.validateLocked()
	m.callReadLocked(f)
}

func (m *RWMutex[T]) validateLocked() {
	m.baseValidateLocked()

	m.initOnce.Do(m.initLocked)
}

func (m *RWMutex[T]) initLocked() {
	if checkTypeCanContainPointers(reflect.TypeOf(m.value)) {
		m.errWrap.err = errContainPointers
		panic(m.errWrap)
	}
	m.initialized = true
}

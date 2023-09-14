package safemutex

import "sync"

type mutexBase[T any] struct {
	m       sync.Mutex
	value   T
	errWrap errWrap
}

func (m *mutexBase[T]) baseValidateLocked() {
	if m.errWrap.err != nil {
		panic(m.errWrap)
	}
}

func (m *mutexBase[T]) callLocked(f ReadWriteCallback[T]) {
	hasPanic := true

	defer func() {
		if hasPanic {
			m.errWrap = errWrap{ErrPoisoned}
		}
	}()

	m.value = f(m.value)
	hasPanic = false

}

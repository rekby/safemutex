//go:build go1.19
// +build go1.19

package safemutex

// TryLockWithResult call f.
// f can change internal state by return new state value. The function return fResult result from the f.
// It is a good way for return state snapshot or some part of them to external code.
// If ok is true:
// first result contains fResult from f.
//
// If ok is false:
// It is mean the f was not called, fResult contains zero value
func TryLockWithResult[M mutexWithTryLock[TState], TState any, TResult any](m M, f func(synced TState) (newSynced TState, fResult TResult)) (fResult TResult, ok bool) {
	var res TResult
	locked := m.TryLock(func(synced TState) TState {
		synced, res = f(synced)
		return synced
	})
	return res, locked
}

type mutexWithTryLock[T any] interface {
	*Mutex[T] | *RWMutex[T] | *RWMutexWithPointers[T]
	TryLock(f ReadWriteCallback[T]) bool
}

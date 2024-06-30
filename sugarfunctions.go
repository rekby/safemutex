package safemutex

// LockWithResult call f.
// f can change internal state by return new state value. The function return second result from the f.
// It is a good way for return state snapshot or some part of them to external code.
func LockWithResult[M mutexWithLock[State], State any, TResult any](m M, f func(synced State) (State, TResult)) TResult {
	var res TResult

	m.Lock(func(synced State) State {
		synced, res = f(synced)
		return synced
	})

	return res
}

type mutexWithLock[T any] interface {
	*Mutex[T] | *RWMutex[T] | *RWMutexWithPointers[T]
	Lock(f ReadWriteCallback[T])
}

// RLockWithResult call f. A result of f returned as the result of RLockWithResult
// It is a good way for return state snapshot or some part of them to external code.
func RLockWithResult[M mutexWithRLock[State], State any, TResult any](m M, f func(synced State) TResult) TResult {
	var res TResult

	m.RLock(func(synced State) {
		res = f(synced)
	})
	return res
}

type mutexWithRLock[T any] interface {
	*RWMutex[T] | *RWMutexWithPointers[T]
	RLock(f ReadCallback[T])
}

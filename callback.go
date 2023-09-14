package safe_mutex

type ReadWriteCallback[T any] func(value T) (newValue T)

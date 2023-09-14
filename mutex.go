package safe_mutex

import (
	"reflect"
	"sync"
)

type Mutex[T any] struct {
	m           sync.Mutex
	value       T
	options     MutexOptions
	initialized bool
	errWrap     errWrap
}

func New[T any](value T) Mutex[T] {
	return NewWithOptions(value, MutexOptions{})
}

func NewWithOptions[T any](value T, options MutexOptions) Mutex[T] {
	res := Mutex[T]{
		value:   value,
		options: options,
	}

	res.validateLocked()

	//goland:noinspection GoVetCopyLock
	return res
}

func (m *Mutex[T]) Lock(f ReadWriteCallback[T]) {
	m.m.Lock()
	defer m.m.Unlock()

	m.callLocked(f)
}

func (m *Mutex[T]) callLocked(f ReadWriteCallback[T]) {
	m.validateLocked()

	hasPanic := true

	defer func() {
		if hasPanic && !m.options.AllowPoisoned {
			m.errWrap = errWrap{ErrMutexPoisoned}
		}
	}()

	m.value = f(m.value)
	hasPanic = false

}

func (m *Mutex[T]) validateLocked() {
	if m.errWrap.err != nil {
		panic(m.errWrap)
	}

	if m.initialized {
		return
	}

	m.initialized = true

	if !m.options.AllowPointers {
		if checkTypeCanContainPointers(reflect.TypeOf(m.value)) {
			panic(errContainPointers)
		}
	}
}

type MutexOptions struct {
	AllowPointers bool
	AllowPoisoned bool
}

// checkTypeCanContainPointers check the value for potential contain pointers
// return true only of value guaranteed without any pointers and false in other cases (has pointers or unknown)
func checkTypeCanContainPointers(t reflectType) bool {
	if t == nil {
		return true
	}
	switch t.Kind() {
	case
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Bool, reflect.Complex64, reflect.Complex128, reflect.Float32, reflect.Float64,
		reflect.String:
		return false
	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			structField := t.Field(i)
			if checkTypeCanContainPointers(structField.Type) {
				return true
			}
		}
		return false
	case reflect.Array:
		return checkTypeCanContainPointers(t.Elem())
	case reflect.Pointer, reflect.UnsafePointer, reflect.Slice, reflect.Map, reflect.Chan, reflect.Interface,
		reflect.Func:
		return true
	default:
		return true
	}
}

type reflectType interface {
	Kind() reflect.Kind
	NumField() int
	Field(i int) reflect.StructField
	Elem() reflect.Type
}

package safe_mutex

import "errors"

var errContainPointers = errors.New("safe mutex: value type possible to contain pointers")
var ErrMutexPoisoned = errors.New("safe mutex: mutex poisoned")

// errWrap need for deny direct compare with returned errors
type errWrap struct {
	err error
}

func (e errWrap) Error() string {
	return e.err.Error()
}

func (e errWrap) Unwrap() error {
	return e.err
}

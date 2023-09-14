package safe_mutex

import (
	"errors"
	"testing"
)

func TestErrWrap_Error(t *testing.T) {
	test := errors.New("test")
	wrap := errWrap{err: test}
	if wrap.Error() != test.Error() {
		t.Fatal(test.Error())
	}
}

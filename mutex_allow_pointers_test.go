package safemutex

import "testing"

func TestNewWithPointers(t *testing.T) {
	val := 1
	t.Run("WithPointers", func(t *testing.T) {
		v := NewWithPointers(struct{ v *int }{&val})
		if v.value.v != &val {
			t.Fatal()
		}
	})
	t.Run("WithoutPointers", func(t *testing.T) {
		v := NewWithPointers(struct{ v int }{val})
		if v.value.v != val {
			t.Fatal()
		}
	})
}

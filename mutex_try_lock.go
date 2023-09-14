//go:build go1.19
// +build go1.19

package safe_mutex

func (m *Mutex[T]) TryLock(f ReadWriteCallback[T]) bool {
	locked := m.m.TryLock()
	if !locked {
		return false
	}
	defer m.m.Unlock()

	m.callLocked(f)
	return true
}

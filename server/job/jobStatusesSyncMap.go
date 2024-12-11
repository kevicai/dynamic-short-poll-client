package job

import (
	"sync"
)

// Adds type safety to the sync.Map
type TypedSyncMap struct {
	data sync.Map
}

func (m *TypedSyncMap) Store(key string, value JobStatus) {
	m.data.Store(key, value)
}

func (m *TypedSyncMap) Load(key string) (status JobStatus, ok bool) {
	val, ok := m.data.Load(key)
	status = val.(JobStatus)
	return
}

func (m *TypedSyncMap) Delete(key string) {
	m.data.Delete(key)
}

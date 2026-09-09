package processor

import (
	"runtime"
	"sync"
)

// MemoryStats contains heap allocation details.
type MemoryStats struct {
	AllocBytes   uint64
	TotalAlloc   uint64
	SysBytes     uint64
	NumGC        uint32
}

// MemoryManager monitors memory and triggers garbage collection every N items.
//
// Requirements: 19.1, 19.2, 19.7, 22.5
type MemoryManager struct {
	mu         sync.Mutex
	gcInterval int
	lastGCCount int
}

// NewMemoryManager creates a new MemoryManager with specified interval.
func NewMemoryManager(gcInterval int) *MemoryManager {
	if gcInterval <= 0 {
		gcInterval = 10
	}
	return &MemoryManager{
		gcInterval: gcInterval,
	}
}

// CheckAndGC triggers runtime.GC() if processedCount reached the interval.
func (mm *MemoryManager) CheckAndGC(processedCount int) bool {
	mm.mu.Lock()
	defer mm.mu.Unlock()

	if processedCount > 0 && (processedCount-mm.lastGCCount) >= mm.gcInterval {
		runtime.GC()
		mm.lastGCCount = processedCount
		return true
	}
	return false
}

// Stats returns current runtime memory statistics.
func (mm *MemoryManager) Stats() MemoryStats {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return MemoryStats{
		AllocBytes: m.Alloc,
		TotalAlloc: m.TotalAlloc,
		SysBytes:   m.Sys,
		NumGC:      m.NumGC,
	}
}

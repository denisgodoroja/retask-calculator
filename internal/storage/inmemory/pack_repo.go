package inmemory

import (
	"sort"
	"sync"
)

// InMemoryPackRepo implements the storage.PackRepository interface using a thread-safe in-memory slice.
type InMemoryPackRepo struct {
	// mu is a Read-Write mutex to protect the sizes slice from concurrent access.
	mu    sync.RWMutex
	sizes []int
}

// NewInMemoryPackRepo creates a new in-memory repository.
func NewInMemoryPackRepo() *InMemoryPackRepo {
	return &InMemoryPackRepo{
		// Start with a default, sorted list
		sizes: []int{250, 500, 1000, 2000, 5000},
	}
}

// FindAll returns a copy of all current pack sizes.
func (r *InMemoryPackRepo) FindAll() ([]int, error) {
	// Read Lock allows multiple concurrent readers.
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Return a copy to prevent the caller from modifying the original slice.
	out := make([]int, len(r.sizes))
	copy(out, r.sizes)

	return out, nil
}

// ReplaceAll replaces all pack sizes with a new list sorted ascending.
func (r *InMemoryPackRepo) ReplaceAll(sizes []int) error {
	// Write Lock blocks all other readers and writers.
	r.mu.Lock()
	defer r.mu.Unlock()

	// Store a copy of the incoming slice
	newSizes := make([]int, len(sizes))
	copy(newSizes, sizes)

	// Sort the sizes to ensure consistency
	sort.Ints(newSizes)

	r.sizes = newSizes

	return nil
}

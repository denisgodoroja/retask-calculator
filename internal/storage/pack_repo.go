package storage

// PackRepository defines the contract for all pack size storage operations.
type PackRepository interface {
	// FindAll returns all current pack sizes, sorted ascending.
	FindAll() ([]int, error)

	// ReplaceAll atomically deletes all existing sizes and inserts the new ones.
	ReplaceAll(sizes []int) error
}

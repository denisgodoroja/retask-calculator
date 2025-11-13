package service

import (
	"denisgodoroja/retask/internal/calculator"
	"denisgodoroja/retask/internal/storage"
)

// PackService holds the core business logic.
type PackService struct {
	repo storage.PackRepository
}

// NewPackService creates a new instance of the PackService.
func NewPackService(r storage.PackRepository) *PackService {
	return &PackService{
		repo: r,
	}
}

// GetPackSizes retrieves the current pack sizes from storage.
func (s *PackService) GetPackSizes() ([]int, error) {
	return s.repo.FindAll()
}

// SetPackSizes persists new pack sizes to storage.
func (s *PackService) SetPackSizes(sizes []int) error {
	return s.repo.ReplaceAll(sizes)
}

// Calculate is the core orchestration logic.
func (s *PackService) Calculate(amount int) (map[int]int, error) {
	sizes, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	packs := calculator.Calculate(amount, sizes)

	return packs, nil
}

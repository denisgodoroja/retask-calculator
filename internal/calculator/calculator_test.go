package calculator

import (
	"reflect"
	"testing"
)

func TestCalculate(t *testing.T) {
	// Standard pack sizes for typical tests
	defaultPacks := []int{250, 500, 1000}

	// A set of fixtures to run the test agains
	fixtures := []struct {
		name      string
		amount    int
		packSizes []int
		expected  map[int]int
	}{
		{
			name:      "Zero amount",
			amount:    0,
			packSizes: defaultPacks,
			expected:  map[int]int{},
		},
		{
			name:      "Amount less than the smallest pack",
			amount:    1,
			packSizes: defaultPacks,
			expected:  map[int]int{250: 1},
		},
		{
			name:      "Exact match smallest pack",
			amount:    250,
			packSizes: defaultPacks,
			expected:  map[int]int{250: 1},
		},
		{
			name:      "Exact match mid pack",
			amount:    500,
			packSizes: defaultPacks,
			expected:  map[int]int{500: 1},
		},
		{
			name:      "Exact match largest pack",
			amount:    1000,
			packSizes: defaultPacks,
			expected:  map[int]int{1000: 1},
		},
		{
			name:      "Just over smallest pack (round up)",
			amount:    251,
			packSizes: defaultPacks,
			expected:  map[int]int{500: 1},
		},
		{
			name:      "Combination of packs (exact)",
			amount:    1750,
			packSizes: defaultPacks,
			expected:  map[int]int{1000: 1, 500: 1, 250: 1},
		},
		{
			name:      "Large amount with remainder",
			amount:    12001,
			packSizes: defaultPacks,
			expected:  map[int]int{1000: 12, 250: 1},
		},
		{
			name:      "No pack sizes provided",
			amount:    100,
			packSizes: []int{},
			expected:  map[int]int{},
		},
		{
			name:      "Single pack size available",
			amount:    3,
			packSizes: []int{5},
			expected:  map[int]int{5: 1},
		},
		{
			name:      "Edge case",
			amount:    500000,
			packSizes: []int{23, 31, 53},
			expected:  map[int]int{23: 2, 31: 7, 53: 9429},
		},
	}

	for _, fixture := range fixtures {
		t.Run(fixture.name, func(t *testing.T) {
			result := Calculate(fixture.amount, fixture.packSizes)
			if !reflect.DeepEqual(result, fixture.expected) {
				t.Errorf("Calculate() = %v, expected %v", result, fixture.expected)
			}
		})
	}
}

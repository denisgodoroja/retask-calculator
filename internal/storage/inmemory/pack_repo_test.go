package inmemory

import (
	"reflect"
	"testing"
)

// TestInMemoryPackRepo_FindAll tests the default state and copy behavior.
func TestInMemoryPackRepo_FindAll(t *testing.T) {
	repo := NewInMemoryPackRepo()

	// 1. Test default values
	sizes, err := repo.FindAll()
	if err != nil {
		t.Fatalf("FindAll() returned an unexpected error: %v", err)
	}
	expected := []int{250, 500, 1000, 2000, 5000}
	if !reflect.DeepEqual(sizes, expected) {
		t.Errorf("FindAll() got = %v, want %v", sizes, expected)
	}

	// 2. Test that a copy is returned
	sizes[0] = 9999
	sizes, err = repo.FindAll()
	if err != nil {
		t.Fatalf("FindAll() returned an unexpected error: %v", err)
	}
	// The internal slice should be unchanged
	if !reflect.DeepEqual(sizes, expected) {
		t.Errorf("FindAll() was modified by external slice change. got = %v, want %v", sizes, expected)
	}
}

// TestInMemoryPackRepo_ReplaceAll tests replacing and sorting.
func TestInMemoryPackRepo_ReplaceAll(t *testing.T) {
	repo := NewInMemoryPackRepo()

	tests := []struct {
		name    string
		input   []int
		want    []int
		wantErr bool
	}{
		{
			name:    "Replace with new sorted list",
			input:   []int{10, 20, 30},
			want:    []int{10, 20, 30},
			wantErr: false,
		},
		{
			name:    "Replace with unsorted list",
			input:   []int{20, 50, 10},
			want:    []int{10, 20, 50}, // Should be sorted
			wantErr: false,
		},
		{
			name:    "Replace with empty list",
			input:   []int{},
			want:    []int{},
			wantErr: false,
		},
		{
			name:    "Replace with nil list",
			input:   nil,
			want:    []int{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := repo.ReplaceAll(tt.input); (err != nil) != tt.wantErr {
				t.Errorf("ReplaceAll() error = %v, wantErr %v", err, tt.wantErr)
			}

			got, err := repo.FindAll()
			if err != nil {
				t.Fatalf("FindAll() returned an unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindAll() got = %v, want %v", got, tt.want)
			}
		})
	}
}

package service

import (
	"errors"
	"reflect"
	"testing"

	"denisgodoroja/retask/internal/storage"
)

// mockPackRepository is a mock implementation of the storage.PackRepository interface.
type mockPackRepository struct {
	// FindAll return:
	findAllSizes []int
	findAllErr   error

	// ReplaceAll return:
	replaceAllErr error

	// To check what was passed in
	replaceAllCalledWith []int
}

func (m *mockPackRepository) FindAll() ([]int, error) {
	return m.findAllSizes, m.findAllErr
}

func (m *mockPackRepository) ReplaceAll(sizes []int) error {
	m.replaceAllCalledWith = sizes
	return m.replaceAllErr
}

// TestPackService_GetPackSizes tests the service layer's GetPackSizes.
func TestPackService_GetPackSizes(t *testing.T) {
	errTest := errors.New("some error")

	tests := []struct {
		name    string
		mock    storage.PackRepository
		want    []int
		wantErr bool
	}{
		{
			name: "Successful fetch",
			mock: &mockPackRepository{
				findAllSizes: []int{100, 200},
				findAllErr:   nil,
			},
			want:    []int{100, 200},
			wantErr: false,
		},
		{
			name: "Error from repository",
			mock: &mockPackRepository{
				findAllErr: errTest,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewPackService(tt.mock)
			got, err := s.GetPackSizes()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPackSizes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPackSizes() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestPackService_SetPackSizes tests the service layer's SetPackSizes.
func TestPackService_SetPackSizes(t *testing.T) {
	errTest := errors.New("some error")

	tests := []struct {
		name    string
		input   []int
		mock    *mockPackRepository
		wantErr bool
	}{
		{
			name:    "Successful set",
			input:   []int{100, 200},
			mock:    &mockPackRepository{},
			wantErr: false,
		},
		{
			name:    "Error from repository",
			input:   []int{100, 200},
			mock:    &mockPackRepository{replaceAllErr: errTest},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewPackService(tt.mock)
			err := s.SetPackSizes(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetPackSizes() error = %v, wantErr %v", err, tt.wantErr)
			}
			// Check that the repo was called with the correct data
			if !tt.wantErr && !reflect.DeepEqual(tt.mock.replaceAllCalledWith, tt.input) {
				t.Errorf("ReplaceAll() not called with correct args. got = %v, want = %v",
					tt.mock.replaceAllCalledWith, tt.input)
			}
		})
	}
}

// TestPackService_Calculate tests the orchestration logic.
func TestPackService_Calculate(t *testing.T) {
	errTest := errors.New("some error")

	tests := []struct {
		name    string
		amount  int
		mock    storage.PackRepository
		want    map[int]int
		wantErr bool
	}{
		{
			name:   "Successful calculation",
			amount: 750,
			mock: &mockPackRepository{
				findAllSizes: []int{250, 500, 1000},
			},
			want:    map[int]int{250: 1, 500: 1},
			wantErr: false,
		},
		{
			name:   "Repo error on FindAll",
			amount: 750,
			mock: &mockPackRepository{
				findAllErr: errTest,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:   "Successful calculation for 1 item",
			amount: 1,
			mock: &mockPackRepository{
				findAllSizes: []int{250, 500},
			},
			want:    map[int]int{250: 1},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewPackService(tt.mock)
			got, err := s.Calculate(tt.amount)
			if (err != nil) != tt.wantErr {
				t.Errorf("Calculate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Calculate() got = %v, want %v", got, tt.want)
			}
		})
	}
}

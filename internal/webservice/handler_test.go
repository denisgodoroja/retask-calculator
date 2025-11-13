package webservice

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"denisgodoroja/retask/internal/service"
	"denisgodoroja/retask/internal/storage"
)

// -- This is a mock *repository* --
type mockPackRepository struct {
	storage.PackRepository // Embed the interface for good practice
	FindAllFunc            func() ([]int, error)
	ReplaceAllFunc         func(sizes []int) error
}

func (m *mockPackRepository) FindAll() ([]int, error)      { return m.FindAllFunc() }
func (m *mockPackRepository) ReplaceAll(sizes []int) error { return m.ReplaceAllFunc(sizes) }

// setupTest creates a Handler with a mock service for testing.
func setupTest() (*Handler, *mockPackRepository) {
	// Create the real service with the mock repo
	mockRepo := &mockPackRepository{}
	realService := service.NewPackService(mockRepo)

	// Create the real handler with the real service
	handler := NewHandler(realService)

	// Return the handler *and* the mock repo so we can control it
	return handler, mockRepo
}

func TestHandler_HandleGetPackSizes(t *testing.T) {
	handler, mockRepo := setupTest()

	// Case 1: Success
	t.Run("Success", func(t *testing.T) {
		mockRepo.FindAllFunc = func() ([]int, error) {
			return []int{100, 200}, nil
		}

		req := httptest.NewRequest(http.MethodGet, "/pack/get-sizes", nil)
		rr := httptest.NewRecorder()
		handler.HandleGetPackSizes(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("wrong status. got %d, want %d", rr.Code, http.StatusOK)
		}

		wantBody := `{"sizes":[100,200]}`
		if rr.Body.String() != wantBody {
			t.Errorf("wrong body. got %q, want %q", rr.Body.String(), wantBody)
		}
	})

	// Case 2: Service Error
	t.Run("Service Error", func(t *testing.T) {
		mockRepo.FindAllFunc = func() ([]int, error) {
			return nil, errors.New("db broke")
		}

		req := httptest.NewRequest(http.MethodGet, "/pack/get-sizes", nil)
		rr := httptest.NewRecorder()
		handler.HandleGetPackSizes(rr, req)

		if rr.Code != http.StatusInternalServerError {
			t.Errorf("wrong status. got %d, want %d", rr.Code, http.StatusInternalServerError)
		}

		wantBody := `{"error":"db broke"}`
		if rr.Body.String() != wantBody {
			t.Errorf("wrong body. got %q, want %q", rr.Body.String(), wantBody)
		}
	})

	// Case 3: Bad Method
	t.Run("Bad Method", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/pack/get-sizes", nil)
		rr := httptest.NewRecorder()
		handler.HandleGetPackSizes(rr, req)

		if rr.Code != http.StatusMethodNotAllowed {
			t.Errorf("wrong status. got %d, want %d", rr.Code, http.StatusMethodNotAllowed)
		}
	})
}

func TestHandler_HandleSetPackSizes(t *testing.T) {
	t.Parallel()
	handler, mockRepo := setupTest()

	// Case 1: Success
	t.Run("Success", func(t *testing.T) {
		mockRepo.ReplaceAllFunc = func(sizes []int) error {
			if !reflect.DeepEqual(sizes, []int{10, 20}) {
				t.Error("ReplaceAll not called with correct args")
			}
			return nil
		}

		body := bytes.NewBufferString(`{"sizes":[10,20]}`)
		req := httptest.NewRequest(http.MethodPost, "/pack/set-sizes", body)
		rr := httptest.NewRecorder()
		handler.HandleSetPackSizes(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("wrong status. got %d, want %d", rr.Code, http.StatusOK)
		}
	})

	// Case 2: Bad JSON
	t.Run("Bad JSON", func(t *testing.T) {
		body := bytes.NewBufferString(`{"sizes":[10,20`) // Malformed
		req := httptest.NewRequest(http.MethodPost, "/pack/set-sizes", body)
		rr := httptest.NewRecorder()
		handler.HandleSetPackSizes(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("wrong status. got %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	// Case 3: Service Error
	t.Run("Service Error", func(t *testing.T) {
		mockRepo.ReplaceAllFunc = func(sizes []int) error {
			return errors.New("db write failed")
		}

		body := bytes.NewBufferString(`{"sizes":[10,20]}`)
		req := httptest.NewRequest(http.MethodPost, "/pack/set-sizes", body)
		rr := httptest.NewRecorder()
		handler.HandleSetPackSizes(rr, req)

		if rr.Code != http.StatusInternalServerError {
			t.Errorf("wrong status. got %d, want %d", rr.Code, http.StatusInternalServerError)
		}
	})
}

func TestHandler_HandleCalculateOrder(t *testing.T) {
	t.Parallel()
	handler, mockRepo := setupTest()

	// Case 1: Success
	t.Run("Success", func(t *testing.T) {
		// This handler calls the service, which calls the repo AND the calculator.
		// We only need to mock the repo part.
		mockRepo.FindAllFunc = func() ([]int, error) {
			return []int{250, 500}, nil // Calculator will use these
		}

		body := bytes.NewBufferString(`{"amount":300}`)
		req := httptest.NewRequest(http.MethodPost, "/order", body)
		rr := httptest.NewRecorder()
		handler.HandleCalculate(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("wrong status. got %d, want %d", rr.Code, http.StatusOK)
		}

		// The real calculator will run: 300 with [250, 500] -> {500: 1}
		var resp CalculateResponse
		if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
			t.Fatal("Could not decode response")
		}

		wantPacks := map[int]int{500: 1}
		if !reflect.DeepEqual(resp.Packs, wantPacks) {
			t.Errorf("wrong packs. got %v, want %v", resp.Packs, wantPacks)
		}
	})

	// Case 2: Repo Error
	t.Run("Repo Error", func(t *testing.T) {
		mockRepo.FindAllFunc = func() ([]int, error) {
			return nil, errors.New("repo died")
		}

		body := bytes.NewBufferString(`{"amount":300}`)
		req := httptest.NewRequest(http.MethodPost, "/calculate", body)
		rr := httptest.NewRecorder()
		handler.HandleCalculate(rr, req)

		if rr.Code != http.StatusInternalServerError {
			t.Errorf("wrong status. got %d, want %d", rr.Code, http.StatusInternalServerError)
		}
		if !bytes.Contains(rr.Body.Bytes(), []byte("repo died")) {
			t.Errorf("wrong body. got %s", rr.Body.String())
		}
	})

	// Case 3: Bad JSON
	t.Run("Bad JSON", func(t *testing.T) {
		body := bytes.NewBufferString(`{"amount":`)
		req := httptest.NewRequest(http.MethodPost, "/calculate", body)
		rr := httptest.NewRecorder()
		handler.HandleCalculate(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("wrong status. got %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})
}

package webservice

import (
	"encoding/json"
	"net/http"

	"denisgodoroja/retask/internal/service"
)

type GetSizesResponse struct {
	Sizes []int `json:"sizes"`
}

type SetSizesRequest struct {
	Sizes []int `json:"sizes"`
}

type CalculateRequest struct {
	Amount int `json:"amount"`
}

type CalculateResponse struct {
	Packs map[int]int `json:"packs"`
}

// Handler holds the dependencies for your HTTP handlers,
// which is primarily the PackService.
type Handler struct {
	service *service.PackService
}

// NewHandler creates a new Handler with its dependencies.
func NewHandler(s *service.PackService) *Handler {
	return &Handler{
		service: s,
	}
}

// HandleGetPackSizes handles GET /pack/get-sizes
func (h *Handler) HandleGetPackSizes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondWithError(w, http.StatusMethodNotAllowed, "Invalid method")
		return
	}

	sizes, err := h.service.GetPackSizes()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, GetSizesResponse{Sizes: sizes})
}

// HandleSetPackSizes handles POST /pack/set-sizes
func (h *Handler) HandleSetPackSizes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "Invalid method")
		return
	}

	var req SetSizesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.service.SetPackSizes(req.Sizes); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// HandleCalculate handles POST /calculate
func (h *Handler) HandleCalculate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "Invalid method")
		return
	}

	var req CalculateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	packs, err := h.service.Calculate(req.Amount)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, CalculateResponse{Packs: packs})
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		// This is a server-side error
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "failed to marshal JSON response"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

package webservice

import (
	"net/http"
)

// NewRouter creates and configures a new http.ServeMux (router)
// It wires all application routes to their corresponding handler methods.
func NewRouter(h *Handler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/pack/get-sizes", h.HandleGetPackSizes)
	mux.HandleFunc("/pack/set-sizes", h.HandleSetPackSizes)
	mux.HandleFunc("/calculate", h.HandleCalculate)

	return mux
}

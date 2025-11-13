package webservice

import (
	"net/http"

	"github.com/gorilla/mux"
)

// NewRouter creates and configures a new router.
// It wires all application routes to their corresponding handler methods.
func NewRouter(h *Handler) http.Handler {
	// Create a new router from gorilla/mux
	router := mux.NewRouter()

	router.HandleFunc("/pack/sizes", h.HandleGetPackSizes).Methods(http.MethodGet)
	router.HandleFunc("/pack/sizes", h.HandleSetPackSizes).Methods(http.MethodPost)
	router.HandleFunc("/calculate", h.HandleCalculate).Methods(http.MethodPost)

	return router
}

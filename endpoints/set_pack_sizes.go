package endpoints

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"denisgodoroja/retask/repository"
)

type packSizesPayload struct {
	Sizes []int `json:"sizes"`
}

func SetPackSizesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Make sure the endpoint is called with POST
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Decode the JSON request body.
		var payload packSizesPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			log.Println("Error decoding body:", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Validate the input
		if len(payload.Sizes) == 0 {
			http.Error(w, "Invalid pack sizes", http.StatusBadRequest)
			return
		}

		// Insert the new sizes
		err := repository.InsertPackSizes(db, payload.Sizes)
		if err != nil {
			log.Println("Failed to insert pack sizes:", err)
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}
	}
}

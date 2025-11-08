package endpoints

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"denisgodoroja/retask/calculator"
	"denisgodoroja/retask/repository"
)

type calculatePayload struct {
	Amount int `json:"amount"`
}

func CalculateHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Make sure the endpoint is called with POST
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Decode the JSON request body.
		var payload calculatePayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			log.Println("Error decoding body:", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// Validate the input
		if payload.Amount <= 0 {
			http.Error(w, "Invalid amount", http.StatusBadRequest)
			return
		}

		sizes, err := repository.FindAllPackSizes(db)
		if err != nil {
			log.Println("Failed finding pack sizes:", err)
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}

		packs := calculator.Calculate(payload.Amount, sizes)

		// Generate successful response's body
		output, _ := json.Marshal(struct {
			Packs map[int]int `json:"packs"`
		}{
			Packs: packs,
		})
		fmt.Fprintln(w, string(output))
	}
}

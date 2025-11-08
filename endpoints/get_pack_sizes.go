package endpoints

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"denisgodoroja/retask/repository"
)

func GetPackSizesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Make sure the endpoint is called with GET
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		sizes, err := repository.FindAllPackSizes(db)
		if err != nil {
			log.Println("Failed fetching pack sizes:", err)
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}

		// Generate successful response's body
		output, _ := json.Marshal(struct {
			Sizes []int `json:"sizes"`
		}{
			Sizes: sizes,
		})
		fmt.Fprintln(w, string(output))
	}
}

package webservice

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	// Import the MySQL driver.
	_ "github.com/go-sql-driver/mysql"

	"denisgodoroja/retask/config"
	"denisgodoroja/retask/endpoints"
)

// Start the web service and database connection
func Start(cfg *config.Config) {
	// Start the database connection
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Name)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	defer db.Close()

	// Check the connection
	if err = db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Successfully connected to the database!")

	// Register endpoints
	http.HandleFunc("/calculate", endpoints.CalculateHandler(db))
	http.HandleFunc("/pack/get-sizes", endpoints.GetPackSizesHandler(db))
	http.HandleFunc("/pack/set-sizes", endpoints.SetPackSizesHandler(db))

	// Start listening for incoming HTTP requests
	log.Fatal(http.ListenAndServe(":"+cfg.Server.Port, nil))
}

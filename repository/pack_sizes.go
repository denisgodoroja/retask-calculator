package repository

import (
	"database/sql"
)

// Insert new pack sizes by overwriting the existing ones
func InsertPackSizes(db *sql.DB, sizes []int) error {
	// Start a transaction. This ensures that either all inserts succeed or none do.
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete the existing sizes prior inserting the new ones
	_, err = tx.Exec("DELETE FROM pack_sizes")
	if err != nil {
		return err
	}

	// Prepare the INSERT query
	stmt, err := tx.Prepare("INSERT INTO pack_sizes (size) VALUES (?)")
	if err != nil {
		return err
	}

	// Close the statement when the function returns
	defer stmt.Close()

	// Run the insert query for each size
	for _, size := range sizes {
		_, err := stmt.Exec(size)
		if err != nil {
			return err
		}
	}

	// Commit the transaction if no errors so far
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// Find all registered pack sizes
func FindAllPackSizes(db *sql.DB) ([]int, error) {
	rows, err := db.Query("SELECT size FROM pack_sizes ORDER BY size ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sizes []int

	for rows.Next() {
		var size int
		if err := rows.Scan(&size); err != nil {
			return nil, err
		}
		sizes = append(sizes, size)
	}

	// Check for any errors that occurred during iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return sizes, nil
}

package weatherservice

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

var db *sql.DB

// InitDB initializes the SQLite database connection and creates the table if it doesn't exist.
func InitDB(dataSourceName string) error {
	var err error
	// Open the database, creating it if it doesn't exist.
	db, err = sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return fmt.Errorf("error opening database: %w", err)
	}

	// Check the connection.
	if err = db.Ping(); err != nil {
		return fmt.Errorf("error pinging database: %w", err)
	}

	// SQL statement to create the table.
	// city is the primary key.
	// data stores the JSON blob.
	createTableSQL := `CREATE TABLE IF NOT EXISTS city_data (
        city TEXT PRIMARY KEY,
        data TEXT
    );`

	// Execute the SQL statement.
	_, err = db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("error creating table: %w", err)
	}

	log.Println("Database initialized successfully.")
	return nil
}

// SaveCityData serializes the provided data map to JSON and saves/updates it in the database.
func SaveCityData(city string, data map[string]any) error {
	if db == nil {
		return fmt.Errorf("database is not initialized")
	}

	// Marshal the data map into a JSON byte slice.
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshaling data to JSON: %w", err)
	}

	cityName := strings.ToLower(city)

	// SQL statement to insert or replace data based on the primary key (city).
	// Using REPLACE ensures that if the city already exists, its data is updated.
	insertSQL := `REPLACE INTO city_data (city, data) VALUES (?, ?);`

	// Execute the SQL statement.
	_, err = db.Exec(insertSQL, cityName, string(jsonData))
	if err != nil {
		return fmt.Errorf("error saving data for city %s: %w", city, err)
	}

	log.Printf("Successfully saved data for city: %s", city)
	return nil
}

// GetCityData retrieves the JSON data string for a specific city.
// It returns sql.ErrNoRows if the city is not found.
func GetCityData(city string) (string, error) {
	if db == nil {
		return "", fmt.Errorf("database is not initialized")
	}

	querySQL := `SELECT data FROM city_data WHERE city = ?;`

	cityName := strings.ToLower(city)

	var jsonData string
	err := db.QueryRow(querySQL, cityName).Scan(&jsonData)
	if err != nil {
		// Don't wrap sql.ErrNoRows, return it directly so the caller can check for it.
		if err == sql.ErrNoRows {
			return "", sql.ErrNoRows
		}
		// For other errors, wrap them.
		return "", fmt.Errorf("error querying data for city %s: %w", city, err)
	}

	log.Printf("Retrieved cached data for city: %s", city)
	return jsonData, nil
}

// CloseDB closes the database connection.
func CloseDB() {
	if db != nil {
		db.Close()
		log.Println("Database connection closed.")
	}
}

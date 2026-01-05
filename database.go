package weatherservice

import (
	"bytes"
	"compress/gzip"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
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
	// data stores the compressed JSON blob.
	createTableSQL := `CREATE TABLE IF NOT EXISTS city_data (
        city TEXT PRIMARY KEY,
        data BLOB 
    );`

	// Execute the SQL statement.
	_, err = db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("error creating table: %w", err)
	}

	log.Println("Database initialized successfully (with compressed data blob).")
	return nil
}

// SaveCityData serializes the provided data map to JSON, compresses it,
// and saves/updates it in the database as a BLOB.
func SaveCityData(city string, data map[string]any) error {
	if db == nil {
		return fmt.Errorf("database is not initialized")
	}

	// 1. Marshal the data map into a JSON byte slice.
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshaling data to JSON: %w", err)
	}

	// 2. Compress the JSON data using gzip
	var compressedBuf bytes.Buffer
	gzWriter := gzip.NewWriter(&compressedBuf)
	if _, err := gzWriter.Write(jsonData); err != nil {
		gzwCloseErr := gzWriter.Close() // Try to close even on write error
		log.Printf("Gzip writer close error after write error: %v", gzwCloseErr)
		return fmt.Errorf("error compressing JSON data for city %s: %w", city, err)
	}
	if err := gzWriter.Close(); err != nil { // Close flushes the compressed data
		return fmt.Errorf("error closing gzip writer for city %s: %w", city, err)
	}

	cityName := strings.ToLower(city)

	// 3. Insert or replace compressed BLOB data based on the primary key (city).
	insertSQL := `REPLACE INTO city_data (city, data) VALUES (?, ?);`

	// Execute the SQL statement.
	_, err = db.Exec(insertSQL, cityName, compressedBuf.Bytes())
	if err != nil {
		return fmt.Errorf("error saving compressed data for city %s: %w", city, err)
	}

	log.Printf("Successfully saved compressed data for city: %s", city)
	return nil
}

// GetCityData retrieves the compressed BLOB data for a specific city,
// decompresses it, and returns the original JSON data as a string.
// It returns sql.ErrNoRows if the city is not found.
func GetCityData(city string) (string, error) {
	if db == nil {
		return "", fmt.Errorf("database is not initialized")
	}

	querySQL := `SELECT data FROM city_data WHERE city = ?;`
	cityName := strings.ToLower(city)

	// 1. Query the compressed BLOB data.
	var compressedData []byte
	err := db.QueryRow(querySQL, cityName).Scan(&compressedData)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", sql.ErrNoRows // Not an error, just cache miss
		}
		// For other errors, wrap them.
		return "", fmt.Errorf("error querying compressed data for city %s: %w", city, err)
	}

	// Handle case where blob might be empty/null in DB
	if len(compressedData) == 0 {
		log.Printf("Warning: Found empty data blob for city %s. Returning empty string.", city)
		return "", nil // Or perhaps return sql.ErrNoRows? Decide based on expected behavior.
	}

	// 2. Decompress the data using gzip.
	compressedReader := bytes.NewReader(compressedData)
	gzReader, err := gzip.NewReader(compressedReader)
	if err != nil {
		// Data might not be compressed? Or corrupted.
		log.Printf("Error creating gzip reader for city %s (data might be uncompressed or corrupt): %v", city, err)
		// Depending on requirements, you might return the raw data (if text) or an error.
		// Returning error is safer if compression is expected.
		return "", fmt.Errorf("error reading compressed data for city %s: %w", city, err)
	}
	defer gzReader.Close()

	decompressedJSON, err := io.ReadAll(gzReader)
	if err != nil {
		return "", fmt.Errorf("error decompressing data for city %s: %w", city, err)
	}

	log.Printf("Successfully retrieved and decompressed data for city: %s", city)
	// 3. Return the original JSON data as a string.
	return string(decompressedJSON), nil
}

// GetAllCountries retrieves all unique countries from the database
func GetAllCountries() ([]string, error) {
	if db == nil {
		return nil, fmt.Errorf("database is not initialized")
	}

	querySQL := `SELECT data FROM city_data;`
	rows, err := db.Query(querySQL)
	if err != nil {
		return nil, fmt.Errorf("error querying all cities: %w", err)
	}
	defer rows.Close()

	countriesMap := make(map[string]bool)
	var countries []string

	for rows.Next() {
		var compressedData []byte
		if err := rows.Scan(&compressedData); err != nil {
			log.Printf("Error scanning city data: %v", err)
			continue
		}

		if len(compressedData) == 0 {
			continue
		}

		// Decompress the data
		compressedReader := bytes.NewReader(compressedData)
		gzReader, err := gzip.NewReader(compressedReader)
		if err != nil {
			log.Printf("Error creating gzip reader: %v", err)
			continue
		}

		decompressedJSON, err := io.ReadAll(gzReader)
		gzReader.Close()
		if err != nil {
			log.Printf("Error decompressing data: %v", err)
			continue
		}

		// Parse JSON to get country
		var data map[string]interface{}
		if err := json.Unmarshal(decompressedJSON, &data); err != nil {
			log.Printf("Error unmarshaling data: %v", err)
			continue
		}

		// Extract country from GeneralInfo
		if generalInfo, ok := data["GeneralInfo"].(map[string]interface{}); ok {
			if country, ok := generalInfo["country"].(string); ok && country != "" {
				// Capitalize first letter of each word
				countryLower := strings.ToLower(country)
				words := strings.Fields(countryLower)
				for i, word := range words {
					if len(word) > 0 {
						words[i] = strings.ToUpper(word[:1]) + word[1:]
					}
				}
				country = strings.Join(words, " ")
				if !countriesMap[country] {
					countriesMap[country] = true
					countries = append(countries, country)
				}
			}
		}
	}

	return countries, nil
}

// GetCitiesByCountry retrieves all cities for a specific country from the database
func GetCitiesByCountry(country string) ([]string, error) {
	if db == nil {
		return nil, fmt.Errorf("database is not initialized")
	}

	querySQL := `SELECT city, data FROM city_data;`
	rows, err := db.Query(querySQL)
	if err != nil {
		return nil, fmt.Errorf("error querying all cities: %w", err)
	}
	defer rows.Close()

	var cities []string
	countryLower := strings.ToLower(country)

	for rows.Next() {
		var cityKey string
		var compressedData []byte
		if err := rows.Scan(&cityKey, &compressedData); err != nil {
			log.Printf("Error scanning city data: %v", err)
			continue
		}

		if len(compressedData) == 0 {
			continue
		}

		// Decompress the data
		compressedReader := bytes.NewReader(compressedData)
		gzReader, err := gzip.NewReader(compressedReader)
		if err != nil {
			log.Printf("Error creating gzip reader: %v", err)
			continue
		}

		decompressedJSON, err := io.ReadAll(gzReader)
		gzReader.Close()
		if err != nil {
			log.Printf("Error decompressing data: %v", err)
			continue
		}

		// Parse JSON to get country
		var data map[string]interface{}
		if err := json.Unmarshal(decompressedJSON, &data); err != nil {
			log.Printf("Error unmarshaling data: %v", err)
			continue
		}

		// Extract country from GeneralInfo
		if generalInfo, ok := data["GeneralInfo"].(map[string]interface{}); ok {
			if cityCountry, ok := generalInfo["country"].(string); ok {
				if strings.ToLower(cityCountry) == countryLower {
					// Extract city name from GeneralInfo
					if cityName, ok := generalInfo["city"].(string); ok && cityName != "" {
						cities = append(cities, cityName)
					}
				}
			}
		}
	}

	return cities, nil
}

// CloseDB closes the database connection.
func CloseDB() {
	if db != nil {
		db.Close()
		log.Println("Database connection closed.")
	}
}

package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

// getEnvOrDefault returns the value of the environment variable or the default value if not set
func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func main() {
	// Database configuration from environment variables
	host := getEnvOrDefault("POSTGRES_HOST", "172.16.168.150")
	portStr := getEnvOrDefault("POSTGRES_PORT", "55432")
	dbname := getEnvOrDefault("POSTGRES_DB", "postgres")
	user := getEnvOrDefault("POSTGRES_USER", "postgres")
	password := getEnvOrDefault("POSTGRES_PASSWORD", "postgres")
	tableName := "t_math_sat_temp_classify_oneprep" // The table we want to inspect

	// Parse port to integer
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Invalid port number: %v", err)
	}

	// Connection string
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Connect to database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to open database connection: %v", err)
	}
	defer db.Close()

	// Test connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	fmt.Println("Successfully connected to the PostgreSQL database!")

	// First, let's check if the table exists
	var exists bool
	err = db.QueryRow(`
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = $1
		)`, tableName).Scan(&exists)

	if err != nil {
		log.Fatalf("Error checking if table exists: %v", err)
	}

	if !exists {
		log.Fatalf("Table %s does not exist", tableName)
	}

	fmt.Printf("Table %s exists. Retrieving column information...\n", tableName)

	// Get column information
	rows, err := db.Query(`
		SELECT column_name, data_type, is_nullable, column_default
		FROM information_schema.columns
		WHERE table_schema = 'public'
		AND table_name = $1
		ORDER BY ordinal_position
	`, tableName)

	if err != nil {
		log.Fatalf("Error querying column information: %v", err)
	}
	defer rows.Close()

	fmt.Println("\nColumn Information:")
	fmt.Printf("%-25s %-15s %-10s %-20s\n", "Column Name", "Data Type", "Nullable", "Default Value")
	fmt.Println(strings.Repeat("-", 70))

	for rows.Next() {
		var columnName, dataType, isNullable string
		var columnDefault sql.NullString

		if err := rows.Scan(&columnName, &dataType, &isNullable, &columnDefault); err != nil {
			log.Fatalf("Error scanning row: %v", err)
		}

		defaultVal := "NULL"
		if columnDefault.Valid {
			defaultVal = columnDefault.String
		}

		fmt.Printf("%-25s %-15s %-10s %-20s\n", columnName, dataType, isNullable, defaultVal)
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("Error iterating rows: %v", err)
	}

	// Now, let's retrieve some sample data from the table
	fmt.Println("\nRetrieving sample data...")

	dataRows, err := db.Query(fmt.Sprintf("SELECT * FROM %s LIMIT 5", tableName))
	if err != nil {
		log.Fatalf("Error querying sample data: %v", err)
	}
	defer dataRows.Close()

	// Get column names for the result set
	columns, err := dataRows.Columns()
	if err != nil {
		log.Fatalf("Error getting column names: %v", err)
	}

	// Prepare values for scanning
	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))

	for i := range columns {
		valuePtrs[i] = &values[i]
	}

	// Print column names
	fmt.Println("\nSample Data:")
	for _, col := range columns {
		fmt.Printf("%-25s", col)
	}
	fmt.Println()
	fmt.Println(strings.Repeat("-", len(columns)*25))

	// Print rows
	for dataRows.Next() {
		err := dataRows.Scan(valuePtrs...)
		if err != nil {
			log.Fatalf("Error scanning row: %v", err)
		}

		for _, val := range values {
			// Convert byte slice to string for bytea columns
			if v, ok := val.([]byte); ok {
				fmt.Printf("%-25s", string(v))
			} else if val == nil {
				fmt.Printf("%-25s", "NULL")
			} else {
				fmt.Printf("%-25v", val)
			}
		}
		fmt.Println()
	}

	if err := dataRows.Err(); err != nil {
		log.Fatalf("Error iterating data rows: %v", err)
	}

	fmt.Println("\nDatabase test completed successfully!")
}

package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Sasha125588/event_app/internal/env"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Get DATABASE_URL or construct it from individual components
	databaseURL := env.GetEnvString("DATABASE_URL", "")

	// If DATABASE_URL is not provided, construct it from individual components
	if databaseURL == "" {
		host := env.GetEnvString("DB_HOST", "localhost")
		port := env.GetEnvInt("DB_PORT", 5432)
		user := env.GetEnvString("DB_USER", "postgres")
		password := env.GetEnvString("DB_PASSWORD", "")
		dbname := env.GetEnvString("DB_NAME", "postgres")
		sslmode := env.GetEnvString("DB_SSLMODE", "require")

		databaseURL = fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
			user, password, host, port, dbname, sslmode)
	}

	// Connect using pgx directly (as in documentation example)
	conn, err := pgx.Connect(context.Background(), databaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer conn.Close(context.Background())

	// Example query to test connection
	var version string
	if err := conn.QueryRow(context.Background(), "SELECT version()").Scan(&version); err != nil {
		log.Fatalf("Query failed: %v", err)
	}

	log.Println("✅ Successfully connected to database!")
	log.Println("PostgreSQL version:", version)

	// Test creating a simple table
	createTable := `
		CREATE TABLE IF NOT EXISTS connection_test (
			id SERIAL PRIMARY KEY,
			message TEXT,
			created_at TIMESTAMP DEFAULT NOW()
		)
	`

	if _, err := conn.Exec(context.Background(), createTable); err != nil {
		log.Fatalf("Failed to create test table: %v", err)
	}

	// Insert test data
	insertQuery := `INSERT INTO connection_test (message) VALUES ($1) RETURNING id`
	var testID int
	if err := conn.QueryRow(context.Background(), insertQuery, "Test connection successful!").Scan(&testID); err != nil {
		log.Fatalf("Failed to insert test data: %v", err)
	}

	log.Printf("✅ Inserted test record with ID: %d", testID)

	// Clean up test table
	if _, err := conn.Exec(context.Background(), "DROP TABLE IF EXISTS connection_test"); err != nil {
		log.Printf("Warning: Failed to clean up test table: %v", err)
	}

	log.Println("🎉 Database connection test completed successfully!")
	log.Println("Your Supabase database is ready to use with the task management API.")
}

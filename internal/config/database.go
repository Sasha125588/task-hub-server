package config

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"

	"github.com/Sasha125588/event_app/internal/env"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DatabaseConfig struct {
	DatabaseURL string
}

func NewDatabaseConfig() *DatabaseConfig {
	databaseURL := env.GetEnvString("DATABASE_URL", "")

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

	databaseURL = addPgxParameters(databaseURL)

	return &DatabaseConfig{
		DatabaseURL: databaseURL,
	}
}

func addPgxParameters(databaseURL string) string {
	u, err := url.Parse(databaseURL)
	if err != nil {
		log.Printf("Warning: Could not parse DATABASE_URL, using as-is: %v", err)
		return databaseURL
	}

	query := u.Query()

	query.Set("default_query_exec_mode", "simple_protocol")
	query.Set("statement_cache_capacity", "512")

	u.RawQuery = query.Encode()

	return u.String()
}

func (c *DatabaseConfig) ConnectDB() (*sql.DB, error) {
	db, err := sql.Open("pgx", c.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	log.Println("Successfully connected to database with pgx driver")
	return db, nil
}

func CreateTables(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS tasks (
			id VARCHAR(255) PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			icon_name VARCHAR(100) NOT NULL,
			start_time VARCHAR(50),
			end_time VARCHAR(50),
			due_date TIMESTAMP NOT NULL,
			progress INTEGER DEFAULT 0,
			status VARCHAR(50) NOT NULL CHECK (status IN ('not-started', 'completed', 'in-progress')),
			comments INTEGER DEFAULT 0,
			attachments INTEGER DEFAULT 0,
			links INTEGER DEFAULT 0,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS sub_tasks (
			id VARCHAR(255) PRIMARY KEY,
			task_id VARCHAR(255) NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
			title VARCHAR(255) NOT NULL,
			description TEXT,
			status VARCHAR(50) NOT NULL CHECK (status IN ('not-started', 'completed', 'in-progress')),
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		)`,
		`CREATE INDEX IF NOT EXISTS idx_tasks_due_date ON tasks(due_date)`,
		`CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status)`,
		`CREATE INDEX IF NOT EXISTS idx_sub_tasks_task_id ON sub_tasks(task_id)`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("error creating table: %w", err)
		}
	}

	log.Println("Successfully created tables")
	return nil
}

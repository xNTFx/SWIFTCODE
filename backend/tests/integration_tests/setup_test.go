package tests

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sql.DB

// setup database connection before running tests
func TestMain(m *testing.M) {
	// get absolute path to .env
	envPath, err := filepath.Abs("../../.env")
	if err != nil {
		fmt.Println("Error getting absolute path:", err)
	}

	err = godotenv.Load(envPath)
	if err != nil {
		fmt.Println("Warning: could not load .env file from", envPath, "Error:", err)
	}

	dsn := os.Getenv("POSTGRES_URL")

	if dsn == "" {
		dsn = "postgres://test:test@localhost:5432/testdb?sslmode=disable"
	}

	db, err = sql.Open("postgres", dsn)
	if err != nil {
		panic("failed to connect to test database: " + err.Error())
	}

	code := m.Run()
	db.Close()
	os.Exit(code)
}

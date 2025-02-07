package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

// initdb attempts to connect to the database up to 5 times,
// logging actual errors and returning a reference or failure.
func InitDB(dataSourceName string) (*sql.DB, error) {
	var database *sql.DB
	var err error

	for i := 1; i <= 5; i++ {
		database, err = sql.Open("postgres", dataSourceName)
		if err != nil {
			log.Printf("attempt %d: failed to open database connection: %v", i, err)
			time.Sleep(2 * time.Second)
			continue
		}

		err = database.Ping()
		if err == nil {
			log.Println("successfully connected to the database")
			return database, nil
		}

		log.Printf("attempt %d: database ping failed: %v", i, err)
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("could not connect to database after 5 attempts: %v", err)
}

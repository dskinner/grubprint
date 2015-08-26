package datastore

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	_ "github.com/lib/pq"
)

var (
	db *sql.DB

	connectOnce sync.Once
)

func Connect() {
	connectOnce.Do(func() {
		var err error
		db, err = sql.Open("postgres", dataSourceName())
		if err != nil {
			log.Fatalf("Failed to open db conn: %v\n", err)
		}
		if err = db.Ping(); err != nil {
			log.Fatalf("Failed to ping db conn: %v\n", err)
		}
	})
}

// TODO
func dataSourceName() string {
	user := "postgres"             // os.Getenv("POSTGRES_USER")
	password := "mysecretpassword" // os.Getenv("POSTGRES_PASSWORD")
	database := "food"             // os.Getenv("POSTGRES_DATABASE")

	return fmt.Sprintf("postgres://%s:%s@localhost:5432/%s?sslmode=disable", user, password, database)
}

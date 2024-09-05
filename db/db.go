package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func connectToDb() *sql.DB {
	url := os.Getenv("DB_URL") + "?authToken=" + os.Getenv("DB_TOKEN")

	db, err := sql.Open("libsql", url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s: %s", url, err)
		os.Exit(1)
	}
	return db
}

var dbInstance *sql.DB

func InitDb() {
	dbInstance = connectToDb()
}

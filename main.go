package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Could not load .env file: %v", err)
	}

	db_user := os.Getenv("DB_USER")
	db_root_pw := os.Getenv("DB_ROOT_PW")
	db_host := os.Getenv("DB_HOST")
	db_name := os.Getenv("DB_NAME")
	app_db_user := os.Getenv("APP_DB_USER")
	app_db_pw := os.Getenv("APP_DB_PW")

	connStr := "postgres://" + db_user + ":" + db_root_pw + "@" + db_host + "/postgres?sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error connection to database: %v", err)
	}
	log.Print("DB Connection Made")
	defer db.Close()

	exists, err := databaseExists(db, db_name)
	if err != nil {
		log.Fatalf("Error checking if database %s exists", db_name)
	}

	if exists {
		log.Printf("Database %s exists, exiting...", db_name)
		return
	}

	query := fmt.Sprintf("CREATE ROLE %s WITH LOGIN NOSUPERUSER NOCREATEDB NOCREATEROLE INHERIT NOREPLICATION NOBYPASSRLS CONNECTION LIMIT -1 PASSWORD '%s'", pq.QuoteIdentifier(app_db_user), app_db_pw)
	log.Printf("%s", query)
	_, err = db.Exec(query)
	if err != nil {
		log.Fatalf("Error creating role %s: %v", app_db_user, err)
	}
	log.Printf("%s Role Created", app_db_user)

	create_db_query := fmt.Sprintf("CREATE DATABASE %s WITH OWNER = %s ENCODING = 'UTF8' LOCALE_PROVIDER = 'libc' CONNECTION LIMIT = -1 IS_TEMPLATE = False", pq.QuoteIdentifier(db_name), pq.QuoteIdentifier(app_db_user))
	_, err = db.Exec(create_db_query)
	if err != nil {
		log.Fatalf("Error creating database %s: %v", db_name, err)
	}
	log.Printf("%s database created", db_name)

}

func databaseExists(db *sql.DB, dbName string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1);"
	err := db.QueryRow(query, dbName).Scan(&exists)
	return exists, err
}

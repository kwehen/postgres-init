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

	var db_user string = os.Getenv("DB_USER")
	var db_root_pw string = os.Getenv("DB_ROOT_PW")
	var db_host string = os.Getenv("DB_HOST")
	var db_name string = os.Getenv("DB_NAME")
	var app_db_user string = os.Getenv("APP_DB_USER")

	connStr := "postgres://" + db_user + ":" + db_root_pw + "@" + db_host + "/postgres?sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error connection to database: %v", err)
	}
	log.Print("DB Connection Made")
	defer db.Close()

	query := fmt.Sprintf("CREATE ROLE %s WITH LOGIN NOSUPERUSER NOCREATEDB NOCREATEROLE INHERIT NOREPLICATION NOBYPASSRLS CONNECTION LIMIT -1", pq.QuoteIdentifier(app_db_user))
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

package db

import (
	"os"
	"fmt"
	"log"
	"database/sql"
	
	_ "github.com/jackc/pgx/v5/stdlib"
)

func Start() *sql.DB {
	dbURL:=fmt.Sprintf("postgresql://%s:%s@%s/%s", os.Getenv("DBUSER"), os.Getenv("DBPASSWD"), os.Getenv("DBHOST"), os.Getenv("DBNAME"))
	
	dbCONN, err := sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	
	return dbCONN
}
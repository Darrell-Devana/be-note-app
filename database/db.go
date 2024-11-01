package database

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDatabase() {
	host := os.Getenv("HOST")
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	user := os.Getenv("USER")
	dbname := os.Getenv("DB_NAME")
	pass := os.Getenv("PASSWORD")

	// set up postgres sql to open it.
	psqlSetup := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		host, port, user, dbname, pass)
	db, err := sql.Open("postgres", psqlSetup)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	if errPing := db.Ping(); errPing != nil {
		fmt.Println("Unable to connect to the database:", errPing)
		panic(errPing)
	}
	DB = db
	fmt.Println("Successfully connected to database!")
}

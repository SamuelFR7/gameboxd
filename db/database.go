package db

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var Db *sqlx.DB

func CreateDatabase() error {
    godotenv.Load()
    var (
        dbpassword = os.Getenv("DATABASE_PASSWORD")
        dbuser = os.Getenv("DATABASE_USER")
        dbname = os.Getenv("DATABASE_DB")
        dbhost = os.Getenv("DATABASE_HOST")
	    uri = fmt.Sprintf("user=%s dbname=%s password=%s host=%s sslmode=disable port=5432", dbuser, dbname, dbpassword, dbhost)
    )

    db, err := sqlx.Connect("postgres", uri)
    if err != nil {
        return err
    }

    Db = db

    return nil
}

package db

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var Db *sqlx.DB

func createDatabase() (*sqlx.DB, error ){
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
        return nil, err
    }

    return db, nil 
}

func Init() error {
    db, err := createDatabase()
    if err != nil {
        return err
    }

    schema, err := os.ReadFile("schema.sql")
    if err != nil {
        return err
    }

    db.MustExec(string(schema))

    Db = db
    
    return nil
}

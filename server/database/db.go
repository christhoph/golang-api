package database

import (
	"fmt"
	"log"
	"os"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"

	"github.com/raunofreiberg/kyrene/server"
)

var DB *pg.DB

func Database() *pg.DB {
	port := server.GetEnvWithDefault(os.Getenv("DB_PORT"), "5432")
	addr := fmt.Sprintf("%v:%v", os.Getenv("DB_HOST"), port)

	return pg.Connect(&pg.Options{
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		Database: os.Getenv("DB_NAME"),
		Addr:     addr,
	})
}

func createSchema(db *pg.DB) error {
	fmt.Println("Migrating... => ", os.Getenv("DB_NAME"))

	for _, model := range Models {
		err := db.CreateTable(model, &orm.CreateTableOptions{
			IfNotExists: true,
			FKConstraints: true,
		})

		if err != nil {
			return err
		}
	}

	return nil
}

// Initialize the database connection and set up migrations, if necessary.
func Init() {
	DB = Database()
	err := createSchema(DB)

	if err != nil {
		log.Fatal(err)
	}
}

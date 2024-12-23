package dbmigration

import (
	"database/sql"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
)

func Migrate(conn *sql.DB) {

	log.Println("DB migration started")

	driver, err := postgres.WithInstance(conn, &postgres.Config{})
	if err != nil {
		log.Fatal("error creating postgres instance")
	}

	m, err := migrate.NewWithDatabaseInstance("file://db/migrations", "postgres", driver)
	if err != nil {
		log.Fatal("DB mmigration failed: ", err)
	}

	if err := m.Down(); err != nil {
		log.Println("DB mmigration (down) failed: ", err)
	}

	if err := m.Up(); err != nil {
		log.Println("DB mmigration (up) failed: ", err)
	}

	log.Println("DB migration completed")
}

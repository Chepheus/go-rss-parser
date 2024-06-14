package storage

import (
	"database/sql"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
)

func MigrationsUp(db *sql.DB, migrationSrc string) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		migrationSrc,
		"postgres",
		driver,
	)

	if err != nil {
		log.Fatal(err)
	}

	// if err := m.Down(); err != nil {
	// log.Fatal(err)
	// }

	if err := m.Up(); err != nil && err.Error() != "no change" {
		log.Fatal(err)
	}
}

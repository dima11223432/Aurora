package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var (
		dsn            string
		migrationPath  string
		migrationTable string
	)

	flag.StringVar(&dsn, "db-dsn", "", "postgres dsn")
	flag.StringVar(&migrationPath, "migrations-path", "", "path to migrations")
	flag.StringVar(&migrationTable, "migrations-table", "schema_migrations", "migration table")
	flag.Parse()

	if dsn == "" {
		panic("db-dsn is required")
	}
	if migrationPath == "" {
		panic("migrations-path is required")
	}

	m, err := migrate.New(
		"file://"+migrationPath,
		dsn+"&x-migrations-table="+migrationTable,
	)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")
			return
		}
		panic(err)
	}

	fmt.Println("migrations applied successfully")
}

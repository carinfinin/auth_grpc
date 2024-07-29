package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3" // драйвер sqlite3
	_ "github.com/golang-migrate/migrate/v4/source/file"      // драйвер работы сфайлома
)

func main() {

	var storagePath, migrationPath, migrationsTable string

	flag.StringVar(&storagePath, "storage_path", "", "path to storage")
	flag.StringVar(&migrationPath, "migration_path", "", "path to migration")
	flag.StringVar(&migrationsTable, "migrations_table", "migrations", "name of migrations table")
	flag.Parse()

	if storagePath == "" {
		panic("storage_path is required")
	}
	if migrationPath == "" {
		panic("migration_path is required")
	}

	m, err := migrate.New(
		"file://"+migrationPath,
		fmt.Sprintf("sqlite3://%s?x-migrations-table=%s", storagePath, migrationsTable),
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

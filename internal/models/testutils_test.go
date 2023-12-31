package models

import (
	"database/sql"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

func migrateDB(t *testing.T, db *sql.DB, rollback bool) error {
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return err
	}

	fSrc, err := (&file.File{}).Open("../../migrations")
	if err != nil {
		t.Fatal(err)
	}

	m, err := migrate.NewWithInstance(
		"file",
		fSrc,
		"sqlite3", driver)
	if err != nil {
		return err
	}
	if rollback {
		err = m.Down()
	} else {
		err = m.Up()
	}
	if err != nil {
		return err
	}
	return nil
}

func newTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", "test_binme.db")
	if err != nil {
		t.Fatal(err)
	}

	err = migrateDB(t, db, false)
	if err != nil {
		db.Close()
		t.Fatal(err)
	}

	t.Cleanup(func() {
		defer db.Close()

		err := migrateDB(t, db, true)
		if err != nil {
			db.Close()
			t.Fatal(err)
		}
	})

	return db
}

package sqlite

import (
	"database/sql"
	"embed"
	"fmt"
)

//go:embed migrations/*
var migrations embed.FS

func Migrate(db *sql.DB) {
	row := db.QueryRow(`SELECT version FROM version;`)

	var version int

	err := row.Scan(&version)
	if err != nil {
		createVersionTable(db)
	}

	migrate(db, version)
}

func createVersionTable(db *sql.DB) {
	_, err := db.Exec(`CREATE TABLE version(version INTEGER);`)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`INSERT INTO version(version) VALUES (0);`)
	if err != nil {
		panic(err)
	}
}

func migrate(db *sql.DB, version int) {
	for {
		filename := fmt.Sprintf(`migrations/%04d.sql`, version+1)

		file, err := migrations.ReadFile(filename)
		if err != nil {
			return
		}

		_, err = db.Exec(string(file))
		if err != nil {
			fmt.Printf("Couldn't run migration %04d.sql\n", version+1)
			panic(err)
		}

		version++

		_, err = db.Exec(`UPDATE version SET version = ?;`, version)
		if err != nil {
			panic(err)
		}
	}
}

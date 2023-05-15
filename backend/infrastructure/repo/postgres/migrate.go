package postgres

import (
	"embed"
	"fmt"

	"github.com/charmbracelet/log"
)

//go:embed migrations/*
var migrations embed.FS

func Migrate(db Queryable) {
	row := db.QueryRow(`SELECT version FROM version;`)

	var version int

	err := row.Scan(&version)
	if err != nil {
		log.Info("failed to scan version: ", err, "; creating version table")
		createVersionTable(db)
	}

	migrate(db, version)
}

func createVersionTable(db Queryable) {
	_, err := db.Exec(`CREATE TABLE version(version INTEGER);`)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`INSERT INTO version(version) VALUES (0);`)
	if err != nil {
		panic(err)
	}
}

func migrate(db Queryable, version int) {
	for {
		filename := fmt.Sprintf(`migrations/%04d.sql`, version+1)

		file, err := migrations.ReadFile(filename)
		if err != nil {
			return
		}

		_, err = db.Exec(string(file))
		if err != nil {
			log.Fatal("Couldn't run migration, file: ", fmt.Sprintf("%04d.sql\n", version+1), " error: ", err)
		}

		version++

		_, err = db.Exec(`UPDATE version SET version = $1;`, version)
		if err != nil {
			panic(err)
		}
	}
}

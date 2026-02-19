package postgres

import (
	"embed"
	"fmt"

	"go.uber.org/zap"
)

//go:embed migrations/*
var migrations embed.FS

func Migrate(db Queryable, logger *zap.Logger) {
	row := db.QueryRow(`SELECT version FROM version;`)

	var version int

	err := row.Scan(&version)
	if err != nil {
		logger.Info("version table not found, creating", zap.Error(err))
		createVersionTable(db)
	}

	migrate(db, version, logger)
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

func migrate(db Queryable, version int, logger *zap.Logger) {
	for {
		filename := fmt.Sprintf(`migrations/%04d.sql`, version+1)

		file, err := migrations.ReadFile(filename)
		if err != nil {
			return
		}

		_, err = db.Exec(string(file))
		if err != nil {
			logger.Fatal("migration failed",
				zap.String("file", fmt.Sprintf("%04d.sql", version+1)),
				zap.Error(err),
			)
		}

		version++

		_, err = db.Exec(`UPDATE version SET version = $1;`, version)
		if err != nil {
			panic(err)
		}

		logger.Info("migration applied", zap.Int("version", version))
	}
}

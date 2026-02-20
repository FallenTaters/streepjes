package postgres

import (
	"embed"
	"fmt"

	"go.uber.org/zap"
)

//go:embed migrations/*
var migrations embed.FS

func Migrate(db Queryable, logger *zap.Logger) error {
	row := db.QueryRow(`SELECT version FROM version;`)

	var version int

	err := row.Scan(&version)
	if err != nil {
		logger.Info("version table not found, creating", zap.Error(err))
		if err := createVersionTable(db); err != nil {
			return fmt.Errorf("migrate: create version table: %w", err)
		}
	}

	return migrate(db, version, logger)
}

func createVersionTable(db Queryable) error {
	if _, err := db.Exec(`CREATE TABLE version(version INTEGER);`); err != nil {
		return fmt.Errorf("create table: %w", err)
	}
	if _, err := db.Exec(`INSERT INTO version(version) VALUES (0);`); err != nil {
		return fmt.Errorf("insert version: %w", err)
	}
	return nil
}

func migrate(db Queryable, version int, logger *zap.Logger) error {
	for {
		filename := fmt.Sprintf(`migrations/%04d.sql`, version+1)

		file, err := migrations.ReadFile(filename)
		if err != nil {
			return nil
		}

		if _, err := db.Exec(string(file)); err != nil {
			return fmt.Errorf("migration %04d.sql: %w", version+1, err)
		}

		version++

		if _, err := db.Exec(`UPDATE version SET version = $1;`, version); err != nil {
			return fmt.Errorf("update version to %d: %w", version, err)
		}

		logger.Info("migration applied", zap.Int("version", version))
	}
}

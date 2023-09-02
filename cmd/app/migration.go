package app

import (
	"database/sql"
	"fmt"

	dbConfig "github.com/nash567/GoSentinel/pkg/db/config"
	logModel "github.com/nash567/GoSentinel/pkg/logger/model"
	"github.com/pressly/goose/v3"
)

func Migrate(log logModel.Logger, db *sql.DB, dbCfg *dbConfig.Config, migrationPath string) {
	goose.SetVerbose(dbCfg.Verbose)
	if err := goose.SetDialect(dbCfg.Dialect); err != nil {
		log.Errorf("could not set dialect for sql migrations: %w", err)
	}

	if err := goose.Up(db, migrationPath); err != nil {
		fmt.Println(err)
		log.Errorf("could not execute migrations: %w", err)
	}
}

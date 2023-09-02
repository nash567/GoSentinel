package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/nash567/GoSentinel/pkg/db/config"
)

func NewConnection(cfg *config.Config) (*sql.DB, error) {
	connectionString := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		cfg.UserName,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
	)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %w", err)
	}

	// Converts struct field names from camel case to snake case for mapping to PostgreSQL table columns.
	// db.Mapper = reflectx.NewMapperFunc("db", strcase.ToSnake)
	return db, nil
}

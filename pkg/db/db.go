package db

import (
	"fmt"

	"github.com/ettle/strcase"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/nash567/GoSentinel/pkg/db/config"

	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
)

func NewConnection(cfg *config.Config) (*sqlx.DB, error) {
	connectionString := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=allow",
		cfg.UserName,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
	)
	db, err := sqlx.Connect("pgx", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %w", err)
	}

	// Converts struct field names from camel case to snake case for mapping to PostgreSQL table columns.
	db.Mapper = reflectx.NewMapperFunc("db", strcase.ToSnake)
	return db, nil
}

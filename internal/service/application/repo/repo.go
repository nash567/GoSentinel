package repo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/lib/pq"
	"github.com/nash567/GoSentinel/internal/service/application/model"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) RegisterApplication(ctx context.Context, application model.Application) error {
	// create schema
	_, err := r.db.ExecContext(ctx, fmt.Sprintf(createSchema, pq.QuoteIdentifier(application.ID)))
	if err != nil {
		return fmt.Errorf("create schema: %w", err)
	}

	// insert into application table
	_, err = r.db.ExecContext(ctx,
		createApplication,
		application.ID,
		application.Name,
		application.Email,
		application.Status,
		application.IsVerified,
	)
	if err != nil {
		return fmt.Errorf("insert application: %w", err)
	}

	// create users table
	_, err = r.db.ExecContext(ctx, fmt.Sprintf(createUsersTable, pq.QuoteIdentifier(application.ID)))
	if err != nil {
		return fmt.Errorf("create users: %w", err)
	}

	return nil
}

func (r *Repository) GetApplication(ctx context.Context, filter *model.Filter) (*model.Application, error) {
	var application model.Application
	q, values := buildQuery(filter)

	err := r.db.QueryRowContext(ctx, q, values...).Scan(
		&application.ID,
		&application.Name,
		&application.Email,
		&application.IsVerified,
		&application.Status,
	)
	if err != nil {
		return nil, fmt.Errorf("get workspaces from db: %w", err)
	}
	return &application, nil
}

func (r *Repository) UpdateApplication(ctx context.Context, application *model.UpdateApplication) error {
	result, err := r.db.Exec(UpdateApplication, application.Name, application.ID)
	if err != nil {
		return fmt.Errorf("update file: %w", err)
	}
	updated, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if updated == 0 {
		return sql.ErrNoRows
	}
	return nil
}

package repo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/nash567/GoSentinel/internal/service/auth/model"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) CreateApplicationIdentity(ctx context.Context, identity model.ApplicationIdentity) error {
	// create applications identity
	_, err := r.db.ExecContext(
		ctx,
		createApplicationIdentity,
		identity.ID,
		identity.ApplicationID,
		identity.Secret,
		identity.SecretViewed,
	)

	if err != nil {
		return fmt.Errorf("insert application identity: %w", err)
	}

	return nil
}

func (r *Repository) GetApplicationIdentity(ctx context.Context, applicationID string) (*model.ApplicationIdentity, error) {
	var identity model.ApplicationIdentity

	err := r.db.QueryRowContext(ctx, getApplicationIdentity, applicationID).Scan(
		&identity.ApplicationID,
		&identity.Secret,
		&identity.SecretViewed,
	)
	if err != nil {
		return nil, fmt.Errorf("get application identity db: %w", err)
	}
	return &identity, nil
}

func (r *Repository) UpdateApplicationIdentity(ctx context.Context, applicationID string) error {
	result, err := r.db.Exec(updateApplicationIdentity, applicationID)
	if err != nil {
		return fmt.Errorf("update application identity: %w", err)
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

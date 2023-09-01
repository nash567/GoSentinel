package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
	"github.com/nash567/GoSentinel/internal/service/user/model"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) RegisterUser(ctx context.Context, user model.User, applicationID string) error {
	result, err := r.db.ExecContext(
		ctx,
		fmt.Sprintf(createUser, pq.QuoteIdentifier(applicationID)),
		user.ID,
		user.Name,
		user.Email,
		user.IsActive,
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert user in db failed: %w", err)
	}
	updatedRows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if updatedRows == 0 {
		return sql.ErrNoRows
	}

	return nil
}
func (r *Repository) UpdateUser(ctx context.Context, user model.UpdateUser, applicationID string) error {
	result, err := r.db.ExecContext(
		ctx,
		fmt.Sprintf(updateUser, pq.QuoteIdentifier(applicationID)),
		user.Name,
		user.ID,
	)
	if err != nil {
		return fmt.Errorf("update user in db failed: %w", err)
	}
	updatedRows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if updatedRows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *Repository) DeleteUser(ctx context.Context, userID string, applicationID string) error {
	result, err := r.db.ExecContext(ctx, fmt.Sprintf(deleteUser, pq.QuoteIdentifier(applicationID)), userID)
	if err != nil {
		return fmt.Errorf("delete user from db failed: %w", err)
	}
	updatedRows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if updatedRows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *Repository) GetUser(ctx context.Context, filter model.Filter, applicationID string) (*model.User, error) {
	user := &model.User{}
	q, values := buildFilter(&filter)
	err := r.db.QueryRow(fmt.Sprintf(getUser+q, pq.QuoteIdentifier(applicationID)), values...).Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found in db: %w", err)
		}
		return nil, fmt.Errorf("get user in db failed: %w", err)
	}
	return user, nil
}

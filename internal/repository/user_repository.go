package repository

import (
	"context"
	"database/sql"
	"fmt"

	"suitemedia/internal/models"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	List(ctx context.Context, params models.ListParams) ([]*models.User, int64, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id string) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (id, email, password, first_name, last_name, role, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING created_at, updated_at
	`

	user.ID = uuid.New()

	err := r.db.QueryRowContext(ctx, query,
		user.ID, user.Email, user.Password, user.FirstName, user.LastName, user.Role, user.IsActive,
	).Scan(&user.CreatedAt, &user.UpdatedAt)

	return err
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	query := `
		SELECT id, email, password, first_name, last_name, role, is_active, created_at, updated_at, deleted_at
		FROM users
		WHERE id = $1 AND deleted_at IS NULL
	`

	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.Password, &user.FirstName, &user.LastName,
		&user.Role, &user.IsActive, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}

	return user, err
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, email, password, first_name, last_name, role, is_active, created_at, updated_at, deleted_at
		FROM users
		WHERE email = $1 AND deleted_at IS NULL
	`

	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.Password, &user.FirstName, &user.LastName,
		&user.Role, &user.IsActive, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return user, err
}

func (r *userRepository) List(ctx context.Context, params models.ListParams) ([]*models.User, int64, error) {
	offset := (params.Page - 1) * params.Limit

	// Count total
	var total int64
	countQuery := `SELECT COUNT(*) FROM users WHERE deleted_at IS NULL`
	if params.Search != "" {
		countQuery += ` AND (first_name ILIKE $1 OR last_name ILIKE $1 OR email ILIKE $1)`
		r.db.QueryRowContext(ctx, countQuery, "%"+params.Search+"%").Scan(&total)
	} else {
		r.db.QueryRowContext(ctx, countQuery).Scan(&total)
	}

	// Get users
	query := `
		SELECT id, email, password, first_name, last_name, role, is_active, created_at, updated_at
		FROM users
		WHERE deleted_at IS NULL
	`

	if params.Search != "" {
		query += ` AND (first_name ILIKE $1 OR last_name ILIKE $1 OR email ILIKE $1)`
	}

	query += ` ORDER BY created_at DESC LIMIT $2 OFFSET $3`

	var rows *sql.Rows
	var err error

	if params.Search != "" {
		rows, err = r.db.QueryContext(ctx, query, "%"+params.Search+"%", params.Limit, offset)
	} else {
		rows, err = r.db.QueryContext(ctx, query, params.Limit, offset)
	}

	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	users := make([]*models.User, 0)
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(
			&user.ID, &user.Email, &user.Password, &user.FirstName, &user.LastName,
			&user.Role, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}

	return users, total, nil
}

func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users
		SET first_name = $1, last_name = $2, role = $3, is_active = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $5 AND deleted_at IS NULL
	`

	_, err := r.db.ExecContext(ctx, query,
		user.FirstName, user.LastName, user.Role, user.IsActive, user.ID,
	)

	return err
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	query := `UPDATE users SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1 AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

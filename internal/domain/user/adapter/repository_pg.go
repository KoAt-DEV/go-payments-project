package adapter

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"go-payments-portfolio-project/internal/database"
	"go-payments-portfolio-project/internal/domain/user"
	"go-payments-portfolio-project/internal/metrics"
)

type PostgresRepository struct {
	pool *database.PgPool
}

func NewPostgresRepository(pool *database.PgPool) user.Repository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) Create(ctx context.Context, u *user.User) error {

	u.CreatedAt = time.Now().UTC()
	u.UpdatedAt = u.CreatedAt

	startDBInsert := time.Now()
	query := `
	INSERT INTO users (email, password_hash, role, coffee_count, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id
	`
	err := r.pool.Pool.QueryRow(ctx, query,
		u.Email,
		u.PasswordHash,
		u.Role,
		u.CoffeeCount,
		u.CreatedAt,
		u.UpdatedAt,
	).Scan(&u.ID)

	duration := time.Since(startDBInsert).Seconds()
	metrics.RegisterDBDurationSeconds.Observe(duration)

	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	u := &user.User{}

	query := `
	SELECT id, email, password_hash, role, coffee_count, created_at, updated_at 
	FROM users 
	WHERE email = $1
	`

	err := r.pool.Pool.QueryRow(ctx, query, email).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.Role,
		&u.CoffeeCount, &u.CreatedAt, &u.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pgx.ErrNoRows
		}
		return nil, err
	}

	return u, nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	u := &user.User{}

	query := `
        SELECT id, email, password_hash, role, coffee_count, created_at, updated_at
        FROM users
        WHERE id = $1
    `

	err := r.pool.Pool.QueryRow(ctx, query, id).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.Role, &u.CoffeeCount,
		&u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pgx.ErrNoRows
		}
		return nil, err
	}

	return u, nil
}

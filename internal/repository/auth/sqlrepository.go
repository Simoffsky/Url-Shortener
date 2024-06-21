package auth

import (
	"context"
	"fmt"

	"url-shorter/internal/models"

	"github.com/jackc/pgx/v4/pgxpool"
)

type PostgresUserRepository struct {
	db *pgxpool.Pool
}

func ConnectToDB(connStr string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

func NewPostgresUserRepository(db *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) CreateUser(user *models.User) error {
	if exists, _ := r.Exists(user.Login); exists {
		return models.ErrUserExists
	}

	query := `INSERT INTO users (login, password) VALUES ($1, $2)`
	_, err := r.db.Exec(context.Background(), query, user.Login, user.Password)
	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}
	return nil
}

func (r *PostgresUserRepository) Exists(login string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE login=$1)`
	err := r.db.QueryRow(context.Background(), query, login).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking if user exists: %w", err)
	}
	return exists, nil
}

func (r *PostgresUserRepository) GetUser(login string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT login, password FROM users WHERE login=$1`
	err := r.db.QueryRow(context.Background(), query, login).Scan(&user.Login, &user.Password)
	if err != nil {
		return nil, models.ErrUserNotFound
	}
	return user, nil
}

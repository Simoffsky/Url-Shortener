package repository

import (
	"context"
	"errors"
	"url-shorter/internal/models"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PgLinksRepository struct {
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

func NewPgLinksRepository(db *pgxpool.Pool) *PgLinksRepository {
	return &PgLinksRepository{db: db}
}

func (r *PgLinksRepository) CreateLink(link models.Link) error {
	sql := `INSERT INTO links (url, short_url, expired_at, creator_login) VALUES ($1, $2, $3, $4)`

	_, err := r.db.Exec(context.Background(), sql, link.Url, link.ShortUrl, link.ExpiredAt, link.CreatorLogin)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // 23505 - unique_violation
			return models.ErrLinkAlreadyExists
		}
		return err
	}
	return nil
}

func (r *PgLinksRepository) GetLink(shortUrl string) (*models.Link, error) {
	sql := `SELECT url, short_url, expired_at, creator_login FROM links WHERE short_url = $1`
	row := r.db.QueryRow(context.Background(), sql, shortUrl)

	var link models.Link
	err := row.Scan(&link.Url, &link.ShortUrl, &link.ExpiredAt, &link.CreatorLogin)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, models.ErrLinkNotFound
	} else if err != nil {
		return nil, err
	}

	return &link, nil
}

func (r *PgLinksRepository) EditLink(shortUrl string, editedLink models.Link) error {
	sql := `UPDATE links SET url = $1, expired_at = $2, creator_login = $3 WHERE short_url = $4`
	_, err := r.db.Exec(context.Background(), sql, editedLink.Url, editedLink.ExpiredAt, editedLink.CreatorLogin, shortUrl)
	if err != nil {
		return models.ErrLinkNotFound
	}
	return nil
}

func (r *PgLinksRepository) RemoveLink(shortUrl string) error {
	sql := `DELETE FROM links WHERE short_url = $1`
	_, err := r.db.Exec(context.Background(), sql, shortUrl)
	if err != nil {
		return models.ErrLinkNotFound
	}
	return nil
}

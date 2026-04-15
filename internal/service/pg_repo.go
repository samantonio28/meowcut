package service

import (
	"database/sql"
	"errors"

	"github.com/samantonio28/meowcut/internal/domain"
)

type pgLinkRepo struct {
	db *sql.DB
}

func NewPGLinkRepo(db *sql.DB) *pgLinkRepo {
	return &pgLinkRepo{db: db}
}

func (r *pgLinkRepo) Save(link domain.Link) error {
	query := `INSERT INTO links (short_id, original_url) VALUES ($1, $2)`
	_, err := r.db.Exec(query, link.ShortID, link.OriginalURL)
	if err != nil {
		return err
	}
	return nil
}

func (r *pgLinkRepo) FindByShortID(shortID string) (domain.Link, error) {
	query := `SELECT short_id, original_url FROM links WHERE short_id = $1`
	var link domain.Link
	err := r.db.QueryRow(query, shortID).Scan(&link.ShortID, &link.OriginalURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Link{}, domain.ErrLinkNotFound
		}
		return domain.Link{}, domain.ErrStorage
	}
	return link, nil
}

func (r *pgLinkRepo) FindByOriginalURL(originalURL string) (domain.Link, error) {
	query := `SELECT short_id, original_url FROM links WHERE original_url = $1`
	var link domain.Link
	err := r.db.QueryRow(query, originalURL).Scan(&link.ShortID, &link.OriginalURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Link{}, domain.ErrLinkNotFound
		}
		return domain.Link{}, domain.ErrStorage
	}
	return link, nil
}

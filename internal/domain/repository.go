package domain

type LinkRepository interface {
	Save(link Link) error
	FindByShortID(shortID string) (Link, error)
	FindByOriginalURL(originalURL string) (Link, error)
}
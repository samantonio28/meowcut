package domain

type LinkCutter interface {
	Cut(originalURL string) (string, error)
}
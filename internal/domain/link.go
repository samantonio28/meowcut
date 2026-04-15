package domain

// Link представляет собой связь между оригинальным URL и коротким идентификатором
type Link struct {
	ShortID     string
	OriginalURL string
}

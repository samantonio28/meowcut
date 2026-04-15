package service

import (
	"sync"

	"github.com/samantonio28/meowcut/internal/domain"
)

type nativeLinkRepo struct {
	mu            sync.RWMutex
	byShortID     map[string]*domain.Link
	byOriginalURL map[string]*domain.Link
}

func NewNativeLinkRepo() *nativeLinkRepo {
	return &nativeLinkRepo{
		byShortID:     make(map[string]*domain.Link),
		byOriginalURL: make(map[string]*domain.Link),
	}
}

func (r *nativeLinkRepo) Save(link domain.Link) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.byShortID[link.ShortID]; exists {
		return domain.ErrDuplicate
	}
	if _, exists := r.byOriginalURL[link.OriginalURL]; exists {
		return domain.ErrDuplicate
	}

	stored := link // copy
	r.byShortID[link.ShortID] = &stored
	r.byOriginalURL[link.OriginalURL] = &stored
	return nil
}

func (r *nativeLinkRepo) FindByShortID(shortID string) (domain.Link, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	linkPtr, ok := r.byShortID[shortID]
	if !ok {
		return domain.Link{}, domain.ErrLinkNotFound
	}
	return *linkPtr, nil
}

func (r *nativeLinkRepo) FindByOriginalURL(originalURL string) (domain.Link, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	linkPtr, ok := r.byOriginalURL[originalURL]
	if !ok {
		return domain.Link{}, domain.ErrLinkNotFound
	}
	return *linkPtr, nil
}

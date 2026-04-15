package service

import (
	"testing"

	"github.com/samantonio28/meowcut/internal/domain"
)

func TestNativeLinkRepo_SaveAndFind(t *testing.T) {
	repo := NewNativeLinkRepo()

	link := domain.Link{
		ShortID:     "Abc123_456",
		OriginalURL: "https://example.com",
	}

	err := repo.Save(link)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	found, err := repo.FindByShortID(link.ShortID)
	if err != nil {
		t.Fatalf("FindByShortID failed: %v", err)
	}
	if found.ShortID != link.ShortID || found.OriginalURL != link.OriginalURL {
		t.Errorf("expected %+v, got %+v", link, found)
	}

	found, err = repo.FindByOriginalURL(link.OriginalURL)
	if err != nil {
		t.Fatalf("FindByOriginalURL failed: %v", err)
	}
	if found.ShortID != link.ShortID || found.OriginalURL != link.OriginalURL {
		t.Errorf("expected %+v, got %+v", link, found)
	}
}

func TestNativeLinkRepo_SaveDuplicateShortID(t *testing.T) {
	repo := NewNativeLinkRepo()

	link1 := domain.Link{ShortID: "SameID", OriginalURL: "https://example1.com"}
	link2 := domain.Link{ShortID: "SameID", OriginalURL: "https://example2.com"}

	err := repo.Save(link1)
	if err != nil {
		t.Fatalf("first Save failed: %v", err)
	}

	err = repo.Save(link2)
	if err != domain.ErrDuplicate {
		t.Errorf("expected ErrDuplicate, got %v", err)
	}
}

func TestNativeLinkRepo_SaveDuplicateOriginalURL(t *testing.T) {
	repo := NewNativeLinkRepo()

	link1 := domain.Link{ShortID: "ID1", OriginalURL: "https://same.com"}
	link2 := domain.Link{ShortID: "ID2", OriginalURL: "https://same.com"}

	err := repo.Save(link1)
	if err != nil {
		t.Fatalf("first Save failed: %v", err)
	}

	err = repo.Save(link2)
	if err != domain.ErrDuplicate {
		t.Errorf("expected ErrDuplicate, got %v", err)
	}
}

func TestNativeLinkRepo_FindByShortIDNotFound(t *testing.T) {
	repo := NewNativeLinkRepo()

	_, err := repo.FindByShortID("NonExistent")
	if err != domain.ErrLinkNotFound {
		t.Errorf("expected ErrLinkNotFound, got %v", err)
	}
}

func TestNativeLinkRepo_FindByOriginalURLNotFound(t *testing.T) {
	repo := NewNativeLinkRepo()

	_, err := repo.FindByOriginalURL("https://nonexistent.com")
	if err != domain.ErrLinkNotFound {
		t.Errorf("expected ErrLinkNotFound, got %v", err)
	}
}

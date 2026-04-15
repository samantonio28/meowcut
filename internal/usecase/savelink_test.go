package usecase

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/samantonio28/meowcut/internal/domain"
	"github.com/samantonio28/meowcut/internal/usecase/mocks"
	"go.uber.org/zap"
)

func TestSaveLinkUsecase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockLinkRepository(ctrl)
	mockCutter := mocks.NewMockLinkCutter(ctrl)
	logger := zap.NewNop()

	uc := NewSaveLinkUsecase(mockRepo, mockCutter, logger)

	t.Run("valid URL, new link", func(t *testing.T) {
		originalURL := "https://example.com"
		shortID := "Abc123_456"

		mockRepo.EXPECT().FindByOriginalURL(originalURL).Return(domain.Link{}, domain.ErrLinkNotFound)
		mockCutter.EXPECT().Cut(originalURL).Return(shortID, nil)
		mockRepo.EXPECT().Save(domain.Link{ShortID: shortID, OriginalURL: originalURL}).Return(nil)

		result, err := uc.Execute(originalURL)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result != shortID {
			t.Errorf("expected short ID %s, got %s", shortID, result)
		}
	})

	t.Run("valid URL, already shortened", func(t *testing.T) {
		originalURL := "https://example.com"
		existingShortID := "ExistingID"

		mockRepo.EXPECT().FindByOriginalURL(originalURL).Return(domain.Link{ShortID: existingShortID, OriginalURL: originalURL}, nil)

		result, err := uc.Execute(originalURL)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result != existingShortID {
			t.Errorf("expected short ID %s, got %s", existingShortID, result)
		}
	})

	t.Run("invalid URL", func(t *testing.T) {
		originalURL := "not-a-url"

		// FindByOriginalURL не должен вызываться
		// Cutter не должен вызываться

		_, err := uc.Execute(originalURL)
		if err != domain.ErrInvalidURL {
			t.Errorf("expected ErrInvalidURL, got %v", err)
		}
	})

	t.Run("cutter error", func(t *testing.T) {
		originalURL := "https://example.com"

		mockRepo.EXPECT().FindByOriginalURL(originalURL).Return(domain.Link{}, domain.ErrLinkNotFound)
		mockCutter.EXPECT().Cut(originalURL).Return("", domain.ErrStorage)

		_, err := uc.Execute(originalURL)
		if err != domain.ErrStorage {
			t.Errorf("expected ErrStorage, got %v", err)
		}
	})

	t.Run("save error", func(t *testing.T) {
		originalURL := "https://example.com"
		shortID := "Abc123_456"

		mockRepo.EXPECT().FindByOriginalURL(originalURL).Return(domain.Link{}, domain.ErrLinkNotFound)
		mockCutter.EXPECT().Cut(originalURL).Return(shortID, nil)
		mockRepo.EXPECT().Save(domain.Link{ShortID: shortID, OriginalURL: originalURL}).Return(domain.ErrStorage)

		_, err := uc.Execute(originalURL)
		if err != domain.ErrStorage {
			t.Errorf("expected ErrStorage, got %v", err)
		}
	})
}
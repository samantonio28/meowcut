package usecase

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/samantonio28/meowcut/internal/domain"
	"github.com/samantonio28/meowcut/internal/usecase/mocks"
	"go.uber.org/zap"
)

func TestGetShortenedUsecase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockLinkRepository(ctrl)
	logger := zap.NewNop()

	uc := NewGetShortenedUsecase(mockRepo, logger)

	t.Run("valid short ID, link exists", func(t *testing.T) {
		shortID := "Abc123_456"
		originalURL := "https://example.com"

		mockRepo.EXPECT().FindByShortID(shortID).Return(domain.Link{ShortID: shortID, OriginalURL: originalURL}, nil)

		result, err := uc.Execute(shortID)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result != originalURL {
			t.Errorf("expected original URL %s, got %s", originalURL, result)
		}
	})

	t.Run("invalid short ID length", func(t *testing.T) {
		shortID := "short"
		// FindByShortID не должен вызываться

		_, err := uc.Execute(shortID)
		if err != domain.ErrInvalidShortID {
			t.Errorf("expected ErrInvalidShortID, got %v", err)
		}
	})

	t.Run("invalid short ID characters", func(t *testing.T) {
		shortID := "Abc123-456" // содержит '-'
		// FindByShortID не должен вызываться

		_, err := uc.Execute(shortID)
		if err != domain.ErrInvalidShortID {
			t.Errorf("expected ErrInvalidShortID, got %v", err)
		}
	})

	t.Run("link not found", func(t *testing.T) {
		shortID := "Abc123_456"

		mockRepo.EXPECT().FindByShortID(shortID).Return(domain.Link{}, domain.ErrLinkNotFound)

		_, err := uc.Execute(shortID)
		if err != domain.ErrLinkNotFound {
			t.Errorf("expected ErrLinkNotFound, got %v", err)
		}
	})

	t.Run("repository error", func(t *testing.T) {
		shortID := "Abc123_456"

		mockRepo.EXPECT().FindByShortID(shortID).Return(domain.Link{}, domain.ErrStorage)

		_, err := uc.Execute(shortID)
		if err != domain.ErrStorage {
			t.Errorf("expected ErrStorage, got %v", err)
		}
	})
}
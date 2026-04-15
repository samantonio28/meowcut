package usecase

import (
	"github.com/samantonio28/meowcut/internal/domain"
	"github.com/samantonio28/meowcut/pkg/utils"
	"go.uber.org/zap"
)

type GetShortenedUsecase struct {
	repo   domain.LinkRepository
	logger *zap.Logger
}

func NewGetShortenedUsecase(repo domain.LinkRepository, logger *zap.Logger) *GetShortenedUsecase {
	return &GetShortenedUsecase{
		repo:   repo,
		logger: logger,
	}
}

func (uc *GetShortenedUsecase) Execute(shortID string) (string, error) {
	if !utils.IsValidShortID(shortID) {
		uc.logger.Warn("invalid short ID provided", zap.String("short_id", shortID))
		return "", domain.ErrInvalidShortID
	}

	link, err := uc.repo.FindByShortID(shortID)
	if err != nil {
		uc.logger.Error("failed to find by short ID", zap.Error(err))
		return "", err
	}

	uc.logger.Info("expanded short link", zap.String("short_id", shortID))
	return link.OriginalURL, nil
}
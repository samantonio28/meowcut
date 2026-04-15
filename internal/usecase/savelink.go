package usecase

import (
	"github.com/samantonio28/meowcut/internal/domain"
	"github.com/samantonio28/meowcut/pkg/utils"
	"go.uber.org/zap"
)

type SaveLinkUsecase struct {
	repo   domain.LinkRepository
	cutter domain.LinkCutter
	logger *zap.Logger
}

func NewSaveLinkUsecase(repo domain.LinkRepository, cutter domain.LinkCutter, logger *zap.Logger) *SaveLinkUsecase {
	return &SaveLinkUsecase{
		repo:   repo,
		cutter: cutter,
		logger: logger,
	}
}

func (uc *SaveLinkUsecase) Execute(originalURL string) (string, error) {
	if !utils.IsValidURL(originalURL) {
		uc.logger.Warn("invalid URL provided", zap.String("url", originalURL))
		return "", domain.ErrInvalidURL
	}

	existing, err := uc.repo.FindByOriginalURL(originalURL)
	if err == nil {
		uc.logger.Info("URL already shortened", zap.String("short_id", existing.ShortID))
		return existing.ShortID, nil
	}
	if err != domain.ErrLinkNotFound {
		uc.logger.Error("failed to find by original URL", zap.Error(err))
		return "", err
	}

	shortID, err := uc.cutter.Cut(originalURL)
	if err != nil {
		uc.logger.Error("failed to generate short ID", zap.Error(err))
		return "", err
	}

	link := domain.Link{
		ShortID:     shortID,
		OriginalURL: originalURL,
	}
	if err := uc.repo.Save(link); err != nil {
		uc.logger.Error("failed to save link", zap.Error(err))
		return "", err
	}

	uc.logger.Info("new short link created", zap.String("short_id", shortID))
	return shortID, nil
}

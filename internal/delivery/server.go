package delivery

import (
	"encoding/json"
	"net/http"

	"github.com/samantonio28/meowcut/internal/api"
	"github.com/samantonio28/meowcut/internal/domain"
	"github.com/samantonio28/meowcut/internal/usecase"
	"go.uber.org/zap"
)

type server struct {
	saveUsecase *usecase.SaveLinkUsecase
	getUsecase  *usecase.GetShortenedUsecase
	logger      *zap.Logger
}

func NewServer(saveUsecase *usecase.SaveLinkUsecase, getUsecase *usecase.GetShortenedUsecase, logger *zap.Logger) api.ServerInterface {
	return &server{
		saveUsecase: saveUsecase,
		getUsecase:  getUsecase,
		logger:      logger,
	}
}

func (s *server) PostShorten(w http.ResponseWriter, r *http.Request) {
	var req api.ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.logger.Warn("failed to decode request", zap.Error(err))
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	shortID, err := s.saveUsecase.Execute(req.Url)
	if err != nil {
		s.logger.Error("shorten usecase failed", zap.Error(err))
		switch err {
		case domain.ErrInvalidURL:
			writeError(w, http.StatusBadRequest, err.Error())
		case domain.ErrDuplicate:
			writeError(w, http.StatusConflict, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	resp := api.ShortenResponse{
		ShortUrl: shortID,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (s *server) GetShortId(w http.ResponseWriter, r *http.Request, shortId string) {
	originalURL, err := s.getUsecase.Execute(shortId)
	if err != nil {
		s.logger.Error("expand usecase failed", zap.Error(err))
		switch err {
		case domain.ErrInvalidShortID:
			writeError(w, http.StatusBadRequest, err.Error())
		case domain.ErrLinkNotFound:
			writeError(w, http.StatusNotFound, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	resp := api.ExpandResponse{
		OriginalUrl: originalURL,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func writeError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(api.Error{Message: message})
}
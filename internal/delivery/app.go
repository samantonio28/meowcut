package delivery

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/samantonio28/meowcut/internal/api"
	"github.com/samantonio28/meowcut/internal/domain"
	"github.com/samantonio28/meowcut/internal/service"
	"github.com/samantonio28/meowcut/internal/usecase"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

type App struct {
	router *mux.Router
	logger *zap.Logger
	config *Config
}

type Config struct {
	Server struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname"`
	} `yaml:"database"`
	StorageType domain.StorageType `yaml:"storage_type"`
}

func NewApp(configPath string) (*App, error) {
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	var config Config
	if err := yaml.Unmarshal(configData, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	logger, err := zap.NewProduction()
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	var repo domain.LinkRepository
	switch config.StorageType {
	case domain.StorageTypeNative:
		repo = service.NewNativeLinkRepo()
		logger.Info("using native storage")
	case domain.StorageTypePG:
		connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			config.Database.Host, config.Database.Port, config.Database.User, config.Database.Password, config.Database.DBName)
		db, err := sql.Open("postgres", connStr)
		if err != nil {
			return nil, fmt.Errorf("failed to open database: %w", err)
		}
		if err := db.Ping(); err != nil {
			return nil, fmt.Errorf("failed to ping database: %w", err)
		}
		repo = service.NewPGLinkRepo(db)
		logger.Info("using PostgreSQL storage")
	default:
		return nil, fmt.Errorf("unknown storage type: %s", config.StorageType)
	}

	cutter := service.NewLinkMeowCutter()

	saveUsecase := usecase.NewSaveLinkUsecase(repo, cutter, logger)
	getUsecase := usecase.NewGetShortenedUsecase(repo, logger)

	server := NewServer(saveUsecase, getUsecase, logger)
	handler := api.Handler(server)

	router := mux.NewRouter()
	router.PathPrefix("/").Handler(handler)

	return &App{
		router: router,
		logger: logger,
		config: &config,
	}, nil
}

func (a *App) Run() error {
	addr := fmt.Sprintf("%s:%d", a.config.Server.Host, a.config.Server.Port)
	a.logger.Info("starting server", zap.String("address", addr))
	return http.ListenAndServe(addr, a.router)
}

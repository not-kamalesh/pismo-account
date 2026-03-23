package server

import (
	"encoding/json"
	"log/slog"
	"os"

	"github.com/not-kamalesh/pismo-account/storage"
	"gorm.io/gorm"
)

type AppConfig struct {
	MySQL *storage.Config `json:"mysql"`
}

type Clients struct {
	DB *gorm.DB
}

func InitClients(appConf *AppConfig) (*Clients, error) {

	// Instantiate the Database Connection Object
	db, err := storage.NewGormDB(appConf.MySQL)
	if err != nil {
		slog.Info("Database initialization failed", "error", err)
		return nil, err
	}

	slog.Info("Database client initialized", "db", db)

	// Run Migrations for the tables
	if err := storage.RunAutoMigrations(db); err != nil {
		slog.Info("Database migrations failed", "error", err)
		return nil, err
	}
	slog.Info("Database migrations applied")

	return &Clients{DB: db}, nil
}

func LoadConfig() (*AppConfig, error) {
	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		path = "config/config.json"
	}
	appCfg := &AppConfig{}
	if err := loadJsonFile(path, appCfg); err != nil {
		slog.Info("Failed to load config", "error", err, "path", path)
		return nil, err
	}
	return appCfg, nil
}

func loadJsonFile(filename string, config interface{}) error {
	jsonFile, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonFile, config)
}

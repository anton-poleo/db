package core

import (
	"go.uber.org/zap"
	"my_db/internal/database"
	"my_db/internal/database/storage"
	"my_db/internal/database/storage/engine"
)

func LoadDatabase(log *zap.Logger, engineType string) (*database.Database, error) {
	eng, err := engine.NewEngine(log, engineType)
	if err != nil {
		return nil, err
	}

	store, err := storage.NewStorage(log, eng)
	if err != nil {
		return nil, err
	}

	db, err := database.NewDatabase(log, store)
	if err != nil {
		log.Error("Failed to create in-memory database", zap.Error(err))
		return nil, err
	}

	return db, nil
}

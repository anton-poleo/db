package core

import (
	"go.uber.org/zap"
	"my_db/internal/database/config"
	"my_db/internal/database/network"
)

func BuildTCPServer(log *zap.Logger, conf *config.Config) (*network.TCPServer, error) {
	db, err := LoadDatabase(log, conf.Engine)
	if err != nil {
		log.Error("failed to load database", zap.Error(err))
		return nil, err
	}

	server, err := loadTCPServer(log, db, conf.Network)
	if err != nil {
		log.Error("failed to load TCP server", zap.Error(err))
		return nil, err
	}

	return server, nil
}

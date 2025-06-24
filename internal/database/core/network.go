package core

import (
	"go.uber.org/zap"
	"my_db/internal/database"
	"my_db/internal/database/config"
	"my_db/internal/database/network"
)

func loadTCPServer(log *zap.Logger, db *database.Database, confNet config.Network) (*network.TCPServer, error) {
	//db, err := database.NewDatabase(log, conf.Engine)
	//if err != nil {
	//	log.Error("Failed to create in-memory database", zap.Error(err))
	//	return nil, err
	//}
	server := network.NewTCPServer(
		log,
		db,
		confNet.Address,
		confNet.IdleTimeout,
		confNet.MaxMessageSize,
		confNet.MaxConnections,
	)
	return server, nil
}

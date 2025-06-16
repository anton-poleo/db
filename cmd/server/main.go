package main

import (
	"my_db/internal/database/core"
)

func main() {
	conf := core.MustLoadConfig()
	log := core.MustLoadLogger(conf.Logging)
	defer log.Sync()

	server, err := core.BuildTCPServer(log, conf)
	if err != nil {
		return
	}
	server.Start()
}

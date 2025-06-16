package core

import (
	"io"
	"my_db/internal/database/config"
	"os"
	"strings"
)

const ConfPathENV = "DB_CONFIG_PATH"

func MustLoadConfig() *config.Config {
	confPath := os.Getenv(ConfPathENV)
	var content io.Reader
	if confPath == "" {
		content = strings.NewReader("")
	} else {
		file, err := os.Open(confPath)
		if err != nil {
			panic(err)
		}
		content = file
	}
	conf, err := config.NewConfig(content)
	if err != nil {
		panic(err)
	}
	return conf
}

package config

import (
	"errors"
	"gopkg.in/yaml.v3"
	"io"
	"time"
)

var (
	NilReaderError = errors.New("nil reader")
)

type Network struct {
	Address              string        `yaml:"address"`
	MaxConnections       int           `yaml:"max_connections"`
	MaxMessageSize       int           `yaml:"max_message_size"`
	IdleTimeout          time.Duration `yaml:"idle_timeout"`
	MaxActiveConnections int           `yaml:"max_active_connections"`
}

type Logging struct {
	Level    string `yaml:"level"`
	Output   string `yaml:"output"`
	Encoding string `yaml:"encoding"`
}

type Config struct {
	Engine  string  `yaml:"engine"`
	Network Network `yaml:"network"`
	Logging Logging `yaml:"output"`
}

func defaultConfig() *Config {
	return &Config{
		Engine: "in_memory",
		Network: Network{
			Address:              "127.0.0.1:3223",
			MaxConnections:       100,
			MaxMessageSize:       1024,
			IdleTimeout:          5 * time.Minute,
			MaxActiveConnections: 1000,
		},
		Logging: Logging{
			Level:    "info",
			Output:   "stdout",
			Encoding: "json",
		},
	}
}

func NewConfig(reader io.Reader) (*Config, error) {
	if reader == nil {
		return nil, NilReaderError
	}

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	config := defaultConfig()
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, err
	}

	return config, nil
}

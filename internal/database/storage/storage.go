package storage

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"my_db/internal/database/storage/engine"
)

var UnknownEngine = errors.New("unknown engine type")
var KeyNotFound = errors.New("not found key")

type Engine interface {
	Get(key string) (string, bool)
	Set(key string, value string)
	Delete(key string)
}

type Storage struct {
	engine Engine
	log    *zap.Logger
}

func NewStorage(log *zap.Logger, engineType string) (*Storage, error) {
	switch engineType {
	case engine.InMemoryEngineType:
		eng, err := engine.NewInMemoryEngine(log)
		if err != nil {
			return &Storage{}, err
		}
		return &Storage{eng, log}, nil
	default:
		return &Storage{}, fmt.Errorf("%w: %s", UnknownEngine, engineType)
	}
}

func (s *Storage) Get(key string) (string, error) {
	val, ok := s.engine.Get(key)
	if !ok {
		return val, fmt.Errorf("%w: %s", KeyNotFound, key)
	}
	return val, nil
}

func (s *Storage) Set(key string, value string) error {
	s.engine.Set(key, value)
	return nil
}

func (s *Storage) Delete(key string) error {
	s.engine.Delete(key)
	return nil
}

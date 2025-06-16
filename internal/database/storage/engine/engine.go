package engine

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"my_db/internal/database/storage/engine/in_memory"
)

const InMemoryEngineType = "in_memory"

type Engine interface {
	Get(key string) (string, bool)
	Set(key string, value string)
	Delete(key string)
}

func NewEngine(log *zap.Logger, engineType string) (Engine, error) {
	if log == nil {
		return nil, errors.New("logger is nil")
	}
	switch engineType {
	case InMemoryEngineType:
		eng, err := in_memory.NewEngine(log)
		if err != nil {
			return nil, err
		}
		return eng, nil
	default:
		return nil, fmt.Errorf("unknown engine: %s", engineType)
	}
}

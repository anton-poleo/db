package engine

import (
	"fmt"
	"go.uber.org/zap"
)

var NilLoggerError = fmt.Errorf("logger is nil")

type InMemoryEngine struct {
	data map[string]string
	log  *zap.Logger
}

func NewInMemoryEngine(log *zap.Logger) (*InMemoryEngine, error) {
	if log == nil {
		return &InMemoryEngine{}, NilLoggerError
	}
	log.Info("Creating in-memory engine")
	return &InMemoryEngine{data: make(map[string]string), log: log}, nil
}

func (e *InMemoryEngine) Get(key string) (string, bool) {
	val, ok := e.data[key]
	return val, ok
}

func (e *InMemoryEngine) Set(key string, value string) {
	e.data[key] = value
}

func (e *InMemoryEngine) Delete(key string) {
	delete(e.data, key)
}

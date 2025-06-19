package in_memory

import (
	"fmt"
	"go.uber.org/zap"
)

var NilLoggerError = fmt.Errorf("logger is nil")

type Engine struct {
	data map[string]string
	log  *zap.Logger
}

func NewEngine(log *zap.Logger) (*Engine, error) {
	if log == nil {
		return &Engine{}, NilLoggerError
	}
	return &Engine{data: make(map[string]string), log: log}, nil
}

func (e *Engine) Get(key string) (string, bool) {
	val, ok := e.data[key]
	e.log.Debug("successfully get query", zap.String("key", key), zap.String("value", val))
	return val, ok
}

func (e *Engine) Set(key string, value string) {
	e.data[key] = value
	e.log.Debug("successfully set query", zap.String("key", key), zap.String("value", value))
}

func (e *Engine) Delete(key string) {
	delete(e.data, key)
	e.log.Debug("successfully delete query", zap.String("key", key))
}

package in_memory

import (
	"fmt"
	"go.uber.org/zap"
)

var NilLoggerError = fmt.Errorf("logger is nil")

type Engine struct {
	data *SharedMap
	log  *zap.Logger
}

func NewEngine(log *zap.Logger) (*Engine, error) {
	if log == nil {
		return &Engine{}, NilLoggerError
	}
	data := NewSharedMap(2)
	return &Engine{data: data, log: log}, nil
}

func (e *Engine) Get(key string) (string, bool) {
	val, ok := e.data.Get(key)
	e.log.Debug("successfully get query", zap.String("key", key), zap.String("value", val))
	return val, ok
}

func (e *Engine) Set(key string, value string) {
	e.data.Set(key, value)
	e.log.Debug("successfully set query", zap.String("key", key), zap.String("value", value))
}

func (e *Engine) Delete(key string) {
	e.data.Delete(key)
	e.log.Debug("successfully delete query", zap.String("key", key))
}

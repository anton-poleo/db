package engine

import (
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	"my_db/internal/database/storage/engine/in_memory"
	"testing"
)

func TestNewEngine(t *testing.T) {
	log := zaptest.NewLogger(t)
	t.Parallel()
	tests := []struct {
		Name       string
		Log        *zap.Logger
		EngineType string
		Err        string
	}{
		{"nil logger", nil, "in_memory", "logger is nil"},
		{"unknown engine", log, "my_engine", "unknown engine: my_engine"},
		{"correct new engine", log, "in_memory", ""},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			engine, err := NewEngine(test.Log, test.EngineType)
			if test.Err == "" {
				require.IsType(t, &in_memory.Engine{}, engine)
				require.NoError(t, err)
			} else {
				require.Nil(t, engine)
				require.ErrorContains(t, err, test.Err)
			}

		})
	}
}

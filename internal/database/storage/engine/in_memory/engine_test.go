package in_memory

import (
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	"testing"
)

func TestNewInMemoryEngine(t *testing.T) {
	log := zaptest.NewLogger(t)
	t.Parallel()
	tests := []struct {
		Name string
		Log  *zap.Logger
		Err  error
	}{
		{"nil logger", nil, NilLoggerError},
		{"correct new engine", log, nil},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			engine, err := NewEngine(test.Log)
			require.Equal(t, test.Err, err)
			require.NotNil(t, engine)
		})
	}
}

func TestInMemoryEngine(t *testing.T) {
	engine, _ := NewEngine(zaptest.NewLogger(t))
	engine.Set("key1", "value1")

	val, ok := engine.Get("key1")
	require.True(t, ok)
	require.Equal(t, "value1", val)

	engine.Delete("key1")
	engine.Delete("key2")

	val, ok = engine.Get("key1")
	require.False(t, ok)
	require.Equal(t, "", val)

}

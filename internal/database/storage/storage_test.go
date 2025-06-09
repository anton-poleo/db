package storage

import (
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	"my_db/internal/database/storage/engine"
	"testing"
)

func TestNewStorage(t *testing.T) {
	t.Parallel()
	tests := []struct {
		Name       string
		EngineType string
		Error      error
	}{
		{"InMemory", engine.InMemoryEngineType, nil},
		{"Unknown", "MyEngine", UnknownEngine},
	}
	log := zaptest.NewLogger(t)

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			st, err := NewStorage(log, test.EngineType)
			require.ErrorIs(t, err, test.Error)
			require.IsType(t, &Storage{}, st)
		})
	}
}

func TestStorage_Get(t *testing.T) {
	log := zaptest.NewLogger(t)
	st, _ := NewStorage(log, engine.InMemoryEngineType)

	_ = st.Set("key", "value")

	t.Parallel()
	tests := []struct {
		Name          string
		Key           string
		ExpectedValue string
		ExpectedError error
	}{
		{"exits key", "key", "value", nil},
		{"not exits key", "key2", "", KeyNotFound},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			val, err := st.Get(test.Key)
			require.Equal(t, val, test.ExpectedValue)
			require.ErrorIs(t, err, test.ExpectedError)
		})
	}
}

func TestStorage_Set(t *testing.T) {
	log := zaptest.NewLogger(t)
	st, _ := NewStorage(log, engine.InMemoryEngineType)

	_, err := st.Get("key")
	require.ErrorIs(t, err, KeyNotFound)
	err = st.Set("key", "value")
	val, err := st.Get("key")
	require.NoError(t, err)
	require.Equal(t, "value", val)

	err = st.Set("key", "new value")
	val, err = st.Get("key")
	require.NoError(t, err)
	require.Equal(t, "new value", val)

}

func TestStorage_Delete(t *testing.T) {
	log := zaptest.NewLogger(t)
	st, _ := NewStorage(log, engine.InMemoryEngineType)

	_ = st.Set("key", "value")
	err := st.Delete("key")
	require.NoError(t, err)

	_, err = st.Get("key")
	require.ErrorIs(t, err, KeyNotFound)
}

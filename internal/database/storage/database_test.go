package storage

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	"my_db/internal/database/compute"
	"my_db/internal/database/storage/engine"
	"testing"
)

func TestNewDatabase(t *testing.T) {
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
			db, err := NewDatabase(log, test.EngineType)
			require.ErrorIs(t, err, test.Error)
			require.IsType(t, &Database{}, db)
		})
	}
}

func TestDatabase_Execute_ParseError(t *testing.T) {
	db, _ := NewDatabase(zaptest.NewLogger(t), engine.InMemoryEngineType)
	value, err := db.Execute("SET key value some")

	require.Error(t, err)
	require.Equal(t, ResponseError, value)
}

func TestDatabase_Execute_UnknownCommand(t *testing.T) {
	db, _ := NewDatabase(zaptest.NewLogger(t), engine.InMemoryEngineType)
	value, err := db.Execute("REMOVE key")

	println(errors.Is(err, compute.UnknownCmdError))
	require.Equal(t, ResponseError, value)
}

func TestDatabase_Execute_Get(t *testing.T) {
	t.Parallel()
	db, _ := NewDatabase(zaptest.NewLogger(t), engine.InMemoryEngineType)
	_, _ = db.Execute("SET key value")

	tests := []struct {
		Name          string
		Key           string
		ExpectedValue string
		ExpectedErr   error
	}{
		{"Exist key", "key", "value", nil},
		{"Not exist key", "key2", ResponseError, KeyNotFound},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			cmd := fmt.Sprintf("GET %s", test.Key)
			val, err := db.Execute(cmd)
			require.ErrorIs(t, err, test.ExpectedErr)
			require.Equal(t, test.ExpectedValue, val)
		})
	}

}

func TestDatabase_Execute_Set(t *testing.T) {
	db, _ := NewDatabase(zaptest.NewLogger(t), engine.InMemoryEngineType)
	_, err := db.Execute("GET key")
	require.ErrorIs(t, err, KeyNotFound)

	_, _ = db.Execute("SET key value")
	value, err := db.Execute("GET key")
	require.NoError(t, err)
	require.Equal(t, "value", value)

	_, _ = db.Execute("SET key value2")
	value, err = db.Execute("GET key")
	require.NoError(t, err)
	require.Equal(t, "value2", value)
}

func TestDatabase_Execute_Delete(t *testing.T) {
	db, _ := NewDatabase(zaptest.NewLogger(t), engine.InMemoryEngineType)

	_, _ = db.Execute("SET key value")
	value, _ := db.Execute("GET key")
	require.Equal(t, "value", value)

	_, _ = db.Execute("DEL key")
	value, err := db.Execute("GET key")
	require.ErrorIs(t, err, KeyNotFound)
	require.Equal(t, value, ResponseError)
}

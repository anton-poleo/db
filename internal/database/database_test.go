package database

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	"my_db/internal/database/compute"
	"my_db/internal/database/storage"
	"my_db/internal/database/storage/engine/in_memory"
	"testing"
)

func InitTestStorage(log *zap.Logger) *storage.Storage {
	eng, _ := in_memory.NewEngine(log)
	st, _ := storage.NewStorage(log, eng)
	return st
}

func TestNewDatabase(t *testing.T) {
	log := zaptest.NewLogger(t)
	t.Parallel()
	tests := []struct {
		Name    string
		Storage *storage.Storage
		Error   string
	}{
		{"InMemory", InitTestStorage(log), ""},
		{"Unknown", nil, "storage is nil"},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			db, err := NewDatabase(log, test.Storage)
			require.IsType(t, &Database{}, db)
			if test.Error == "" {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, test.Error)
			}
		})
	}
}

func TestDatabase_Execute_ParseError(t *testing.T) {
	log := zaptest.NewLogger(t)
	db, _ := NewDatabase(log, InitTestStorage(log))
	value, err := db.Execute("SET key value some")

	require.Error(t, err)
	require.Equal(t, ResponseError, value)
}

func TestDatabase_Execute_UnknownCommand(t *testing.T) {
	log := zaptest.NewLogger(t)
	db, _ := NewDatabase(log, InitTestStorage(log))
	value, err := db.Execute("REMOVE key")

	println(errors.Is(err, compute.UnknownCmdError))
	require.Equal(t, ResponseError, value)
}

func TestDatabase_Execute_Get(t *testing.T) {
	t.Parallel()
	log := zaptest.NewLogger(t)
	db, _ := NewDatabase(log, InitTestStorage(log))
	_, _ = db.Execute("SET key value")

	tests := []struct {
		Name          string
		Key           string
		ExpectedValue string
		ExpectedErr   error
	}{
		{"Exist key", "key", "value", nil},
		{"Not exist key", "key2", ResponseError, storage.KeyNotFound},
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
	log := zaptest.NewLogger(t)
	db, _ := NewDatabase(log, InitTestStorage(log))
	_, err := db.Execute("GET key")
	require.ErrorIs(t, err, storage.KeyNotFound)

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
	log := zaptest.NewLogger(t)
	db, _ := NewDatabase(log, InitTestStorage(log))

	_, _ = db.Execute("SET key value")
	value, _ := db.Execute("GET key")
	require.Equal(t, "value", value)

	_, _ = db.Execute("DEL key")
	value, err := db.Execute("GET key")
	require.ErrorIs(t, err, storage.KeyNotFound)
	require.Equal(t, value, ResponseError)
}

package database

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"my_db/internal/database/compute"
)

const (
	ResponseOk    = "Ok"
	ResponseError = "Error"
)

type Storage interface {
	Get(key string) (string, error)
	Set(key string, value string) error
	Delete(key string) error
}

type Database struct {
	storage Storage
	parser  *compute.Parser
	logger  *zap.Logger
}

func NewDatabase(log *zap.Logger, storage Storage) (*Database, error) {
	if storage == nil {
		return nil, errors.New("storage is nil")
	}
	parser := compute.NewParser(log)
	return &Database{storage, parser, log}, nil
}

func (db *Database) Execute(cmd string) (string, error) {
	query, err := db.parser.ParseCommand(cmd)
	if err != nil {
		db.logger.Error("parse command", zap.String("cmd", cmd), zap.Error(err))
		return ResponseError, err
	}

	switch query.Name() {
	case compute.CommandGET:
		return db.executeGet(query.Args())
	case compute.CommandSET:
		return db.executeSet(query.Args())
	case compute.CommandDEL:
		return db.executeDelete(query.Args())
	default:
		return ResponseError, fmt.Errorf("%w: %s", compute.UnknownCmdError, query.Name())
	}
}

func (db *Database) executeGet(args []string) (string, error) {
	val, err := db.storage.Get(args[0])
	if err != nil {
		return ResponseError, err
	}
	return val, nil
}

func (db *Database) executeSet(args []string) (string, error) {
	err := db.storage.Set(args[0], args[1])
	if err != nil {
		return ResponseError, err
	}
	return ResponseOk, nil
}

func (db *Database) executeDelete(args []string) (string, error) {
	err := db.storage.Delete(args[0])
	if err != nil {
		return ResponseError, err
	}
	return ResponseOk, nil
}

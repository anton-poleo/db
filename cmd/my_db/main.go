package main

import (
	"bufio"
	"fmt"
	"go.uber.org/zap"
	"my_db/internal/database"
	"my_db/internal/database/storage"
	"my_db/internal/database/storage/engine"
	"os"
)

func main() {
	config := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.DebugLevel),
		Encoding:    "console",
		OutputPaths: []string{"stdout"},
	}
	log, err := config.Build()
	if err != nil {
		panic(err)
	}
	eng, err := engine.NewEngine(log, engine.InMemoryEngineType)
	if err != nil {
		panic(err)
	}
	store, err := storage.NewStorage(log, eng)
	if err != nil {
		panic(err)
	}
	db, err := database.NewDatabase(log, store)
	if err != nil {
		panic(err)
	}
	log.Info("db started...")

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		cmd := scanner.Text()
		if cmd == "exit" {
			return
		}
		value, err := db.Execute(cmd)
		if err != nil {
			fmt.Printf("error: %s\n", err.Error())
		} else {
			fmt.Println(value)
		}
	}
}

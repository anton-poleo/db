package main

import (
	"bufio"
	"fmt"
	"go.uber.org/zap"
	"my_db/internal/database/storage"
	"os"
)

func main() {
	config := zap.Config{
		Level:    zap.NewAtomicLevelAt(zap.DebugLevel),
		Encoding: "console",
	}
	log, err := config.Build()
	if err != nil {
		panic(err)
	}

	database, err := storage.NewDatabase(log, "in_memory")
	if err != nil {
		panic(err)
	}
	log.Info("database started...")

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		cmd := scanner.Text()
		if cmd == "exit" {
			return
		}
		value, err := database.Execute(cmd)
		if err != nil {
			fmt.Printf("error: %s\n", err.Error())
		} else {
			fmt.Println(value)
		}
	}
}

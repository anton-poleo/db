package main

import (
	"flag"
	"fmt"
	"my_db/internal/database/network"
	"time"
)

func main() {
	address := flag.String("address", "127.0.0.1:3223", "server address")
	idleTimeout := flag.Duration("idleTimeout", 60*time.Second, "idle timeout")
	maxMessageSize := flag.Int("maxMessageSize", 1024*1024, "max message size")
	flag.Parse()

	client := network.NewTCPClient(*address, *idleTimeout, *maxMessageSize)
	err := client.Start()
	if err != nil {
		fmt.Println("Failed to start client", err.Error())
	}
}

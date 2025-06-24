package network

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
)

const ExitCmd = "exit"

type TCPClient struct {
	address        string
	idleTimeout    time.Duration
	maxMessageSize int
}

func NewTCPClient(address string, idleTimeout time.Duration, maxMessageSize int) *TCPClient {
	return &TCPClient{
		address:        address,
		idleTimeout:    idleTimeout,
		maxMessageSize: maxMessageSize,
	}
}

func (c *TCPClient) Start() error {
	fmt.Println("Starting TCP client...")
	conn, err := net.Dial(TCPNetwork, c.address)
	if err != nil {
		fmt.Println("Failed to connect to server", err.Error())
		return err
	}
	println(c.idleTimeout)
	err = conn.SetReadDeadline(time.Now().Add(c.idleTimeout))
	if err != nil {
		fmt.Println("Failed to set read deadline", err.Error())
	}

	fmt.Println("Connected to server", conn.RemoteAddr().String())
	defer conn.Close()

	scanner := bufio.NewScanner(os.Stdin)
	buf := make([]byte, c.maxMessageSize)
	for scanner.Scan() {
		query := scanner.Text()
		if query == ExitCmd {
			return nil
		}
		_, err := conn.Write([]byte(query))
		if err != nil {
			fmt.Println("Failed to write to server", err.Error())
			continue
		}
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Failed to read from server", err.Error())
			continue
		}
		fmt.Println(string(buf[:n]))

		err = conn.SetReadDeadline(time.Now().Add(c.idleTimeout))
		if err != nil {
			fmt.Println("Failed to set read deadline", err.Error())
		}
	}
	return nil
}

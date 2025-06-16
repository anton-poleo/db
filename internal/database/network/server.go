package network

import (
	"go.uber.org/zap"
	"net"
	"time"
)

const TCPNetwork = "tcp"

type TCPServer struct {
	log                  *zap.Logger
	handler              *TCPHandler
	address              string
	maxActiveConnections int
}

func NewTCPServer(
	log *zap.Logger,
	db QueryExecutor,
	address string,
	idleTimeout time.Duration,
	maxMessageSize int,
	maxActiveConnections int,
) *TCPServer {
	handler := NewTCPHandler(log, db, idleTimeout, maxMessageSize)
	return &TCPServer{
		log:                  log,
		handler:              handler,
		address:              address,
		maxActiveConnections: maxActiveConnections,
	}
}

func (s *TCPServer) Start() {
	listener, err := net.Listen(TCPNetwork, s.address)
	if err != nil {
		s.log.Error("error while start listening", zap.Error(err))
		return
	}

	defer listener.Close()

	s.log.Info("server started", zap.String("address", listener.Addr().String()))

	sema := make(chan struct{}, s.maxActiveConnections)

	for {
		conn, err := listener.Accept()
		if err != nil {
			s.log.Error("error while accepting connection", zap.Error(err))
			continue
		}
		select {
		case sema <- struct{}{}:
			s.log.Info("new connection", zap.String("remote", conn.RemoteAddr().String()))
			go s.handler.HandleRequest(conn)
		default:
			s.log.Warn("max connections reached")
			conn.Write([]byte("max connections reached"))
			conn.Close()
		}

	}
}

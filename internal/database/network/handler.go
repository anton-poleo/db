package network

import (
	"fmt"
	"go.uber.org/zap"
	"net"
	"strings"
	"time"
)

type QueryExecutor interface {
	Execute(query string) (string, error)
}

type TCPHandler struct {
	db             QueryExecutor
	log            *zap.Logger
	idleTimeout    time.Duration
	maxMessageSize int
}

func NewTCPHandler(
	log *zap.Logger,
	db QueryExecutor,
	idleTimeout time.Duration,
	maxMessageSize int,
) *TCPHandler {
	return &TCPHandler{
		log:            log,
		db:             db,
		idleTimeout:    idleTimeout,
		maxMessageSize: maxMessageSize,
	}
}

func (h *TCPHandler) HandleRequest(conn net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			h.log.Error("captured panic", zap.Any("error", err))
		}
		if err := conn.Close(); err != nil {
			h.log.Error("failed to close connection", zap.Error(err))
		}
		h.log.Info("connection closed")
	}()

	for {
		deadline := time.Now().Add(h.idleTimeout)
		if err := conn.SetReadDeadline(deadline); err != nil {
			h.log.Error("failed to set read deadline", zap.Error(err))
			return
		}

		buf := make([]byte, h.maxMessageSize)

		n, err := conn.Read(buf)

		if err != nil {
			h.log.Error("failed to read from connection", zap.Error(err))
			h.log.Error("break")
			break
		}
		if n == h.maxMessageSize {
			h.log.Error("message too large", zap.Int("size", n), zap.Int("max_size", h.maxMessageSize))
			if _, err = conn.Write([]byte(fmt.Sprintf("message too large (%d > %d)", n, h.maxMessageSize))); err != nil {
				h.log.Error(
					"failed to write to connection",
					zap.String("address", conn.RemoteAddr().String()),
					zap.Error(err),
				)
			}
			break
		}

		query := strings.TrimSpace(string(buf[:n]))
		if query == "" {
			continue
		}

		value, err := h.db.Execute(string(buf[:n]))

		var msg string
		// todo: получать msg из Execute
		if err != nil {
			msg = fmt.Sprintf("error: %s", err.Error())
		} else {
			msg = fmt.Sprintf(value)
		}

		if _, err := conn.Write([]byte(msg)); err != nil {
			h.log.Error(
				"failed to write to connection",
				zap.String("address", conn.RemoteAddr().String()),
				zap.Error(err),
			)
		}
	}
}

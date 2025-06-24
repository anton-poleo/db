package config

import (
	"errors"
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"strings"
	"testing"
	"time"
)

type BreakingReader struct{}

func (br *BreakingReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("some error")
}

func TestLoadConfig_Error(t *testing.T) {
	tests := []struct {
		name   string
		reader io.Reader
	}{
		{"Nil reader", nil},
		{"Invalid reader message", strings.NewReader("hello")},
		{"Invalid reader", &BreakingReader{}},
	}

	t.Parallel()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := NewConfig(test.reader)
			require.Error(t, err)
		})
	}
}

func TestLoadConfig(t *testing.T) {
	file, _ := os.Open("../../../config/local.yaml")
	defer file.Close()

	cfg, err := NewConfig(file)
	require.NoError(t, err)

	require.Equal(t, cfg.Engine, "in_memory")
	require.Equal(t, cfg.Network.IdleTimeout, 5*time.Minute)
	require.Equal(t, cfg.Network.MaxMessageSize, 1024)
}

package compute

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	"strings"
	"testing"
)

func TestQuery(t *testing.T) {
	name := "GET"
	args := []string{"param1"}

	query := NewQuery(name, args)

	require.Equal(t, name, query.Name())
	require.Equal(t, args, query.Args())
}

func TestParseCommand(t *testing.T) {
	tests := []struct {
		Name string
		Cmd  string
		Args []string
		Err  error
	}{
		{"get ok", "GET", []string{"some_param"}, nil},
		{"set ok", "SET", []string{"some_param", "some_value"}, nil},
		{"empty", "", nil, EmptyCmdError},
		{"unknown cmd", "PUSH 12", nil, UnknownCmdError},
		{"incorrect num args get", "GET key value", nil, IncorrectArgsNumberError},
		{"incorrect num args set", "SET 12", nil, IncorrectArgsNumberError},
		{"incorrect get", "GET 1++2", nil, InvalidArgsError},
	}

	log := zaptest.NewLogger(t)

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			command := fmt.Sprintf("%s %s", tt.Cmd, strings.Join(tt.Args, " "))
			var expectedQuery Query
			if tt.Err == nil {
				expectedQuery = Query{
					name: tt.Cmd,
					args: tt.Args,
				}
			}

			query, err := NewParser(log).ParseCommand(command)

			require.Equal(t, tt.Err, err)
			require.Equal(t, expectedQuery, query)
		})
	}
}

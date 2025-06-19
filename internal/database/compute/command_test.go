package compute

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetCommandArgsNum(t *testing.T) {
	require.Equal(t, CommandGETArgsNum, getCommandArgsNum(CommandGET))
	require.Equal(t, CommandSETArgsNum, getCommandArgsNum(CommandSET))
	require.Equal(t, CommandDELArgsNum, getCommandArgsNum(CommandDEL))
}

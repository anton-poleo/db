package in_memory

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSharedMap(t *testing.T) {
	hashMap := NewSharedMap(3)
	hashMap.Set("foo", "bar")
	hashMap.Set("foo2", "bar2")

	v1, ok := hashMap.Get("foo")
	require.True(t, ok)
	require.Equal(t, "bar", v1)

	v2, ok := hashMap.Get("foo2")
	require.True(t, ok)
	require.Equal(t, "bar2", v2)

	hashMap.Delete("foo")
	v1, ok = hashMap.Get("foo")
	require.False(t, ok)
	//require.Empty(t, v1)
}

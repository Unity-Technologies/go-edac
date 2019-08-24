package edac

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHasMemoryErrors(t *testing.T) {
	clean := mockSysfs(t)
	defer clean()

	ok, err := HasMemoryErrors()
	require.NoError(t, err)
	require.True(t, ok)
}

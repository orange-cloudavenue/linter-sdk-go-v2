package linters

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func testdataDir(t *testing.T) string {
	t.Helper()

	_, testFilename, _, ok := runtime.Caller(1)
	if !ok {
		require.Fail(t, "unable to get current test filename")
	}

	return filepath.Join(filepath.Dir(testFilename), "testdata")
}

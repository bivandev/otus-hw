package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	testdataDir := "testdata"

	t.Run("successful copy entire file", func(t *testing.T) {
		srcPath := filepath.Join(testdataDir, "input.txt")

		tmpFile, err := os.CreateTemp(t.TempDir(), "testfile.txt")
		require.NoError(t, err)

		err = Copy(srcPath, tmpFile.Name(), 0, 0)
		require.NoError(t, err)

		content, err := os.ReadFile(tmpFile.Name())
		require.NoError(t, err)
		srcContent, err := os.ReadFile(srcPath)
		require.NoError(t, err)

		require.Equal(t, srcContent, content)
	})

	t.Run("copy with offset and limit", func(t *testing.T) {
		srcPath := filepath.Join(testdataDir, "input.txt")

		tmpFile, err := os.CreateTemp(t.TempDir(), "testfile.txt")
		require.NoError(t, err)

		err = Copy(srcPath, tmpFile.Name(), 100, 1000)
		require.NoError(t, err)

		content, err := os.ReadFile(tmpFile.Name())
		require.NoError(t, err)

		out, err := os.ReadFile(filepath.Join(testdataDir, "out_offset100_limit1000.txt"))
		require.NoError(t, err)

		require.Equal(t, out, content)
	})

	t.Run("unsupported file type", func(t *testing.T) {
		srcPath := "/dev/null"
		dstPath := t.TempDir()

		err := Copy(srcPath, dstPath, 0, 0)
		require.ErrorIs(t, err, ErrUnsupportedFile)
	})

	t.Run("empty source file", func(t *testing.T) {
		tmpEmptFile, err := os.CreateTemp(t.TempDir(), "testfile.txt")
		require.NoError(t, err)

		tmpFile, err := os.CreateTemp(t.TempDir(), "testfile.txt")
		require.NoError(t, err)

		err = Copy(tmpEmptFile.Name(), tmpFile.Name(), 0, 0)
		require.NoError(t, err)

		content, err := os.ReadFile(tmpFile.Name())
		require.NoError(t, err)
		require.Equal(t, "", string(content))
	})

	t.Run("destination file creation error", func(t *testing.T) {
		srcPath := filepath.Join(testdataDir, "input.txt")
		dstPath := "/invalid/path/dest.txt"

		err := Copy(srcPath, dstPath, 0, 0)
		require.Error(t, err)
	})
}

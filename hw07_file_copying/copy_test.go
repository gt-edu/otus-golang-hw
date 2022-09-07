package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	// Place your code here.

	tmpDir := os.TempDir()
	fromFile, err := os.CreateTemp(tmpDir, "simple-case-*")
	require.NoError(t, err)

	testContentBytes := []byte(("test string"))
	n, err := fromFile.Write(testContentBytes)
	require.NoError(t, err)
	require.Equal(t, len(testContentBytes), n)

	err = fromFile.Close()
	require.NoError(t, err)

	toFilePath := fromFile.Name() + "-target"

	t.Run("validate limit and offset", func(t *testing.T) {
		err = Copy(fromFile.Name(), toFilePath, -1, 0)
		require.EqualError(t, err, ErrOffsetAndLimitShouldBePositive.Error())

		err = Copy(fromFile.Name(), toFilePath, 0, -1)
		require.EqualError(t, err, ErrOffsetAndLimitShouldBePositive.Error())

		err = Copy(fromFile.Name(), toFilePath, int64(len(testContentBytes))+1, 0)
		require.EqualError(t, err, ErrOffsetExceedsFileSize.Error())
	})

	t.Run("validate unsupported file", func(t *testing.T) {
		err = Copy("/dev/random", toFilePath, 0, 0)
		require.EqualError(t, err, ErrUnsupportedFile.Error())

		err = Copy(fromFile.Name(), "/dev/random", 0, 0)
		require.EqualError(t, err, ErrUnsupportedFile.Error())
	})

	t.Run("full copy", func(t *testing.T) {
		err = Copy(fromFile.Name(), toFilePath, 0, 0)
		require.NoError(t, err)

		resultBytes := readFileBytes(t, toFilePath)
		require.Equal(t, testContentBytes, resultBytes)
	})

	t.Run("no limit with offset", func(t *testing.T) {
		testOffset := 5
		err = Copy(fromFile.Name(), toFilePath, int64(testOffset), 0)
		require.NoError(t, err)

		resultBytes := readFileBytes(t, toFilePath)
		require.Equal(t, testContentBytes[testOffset:], resultBytes)
	})

	t.Run("no offset with limit", func(t *testing.T) {
		testLimit := 5
		err = Copy(fromFile.Name(), toFilePath, 0, int64(testLimit))
		require.NoError(t, err)

		resultBytes := readFileBytes(t, toFilePath)
		require.Equal(t, testContentBytes[:testLimit], resultBytes)
	})

	t.Run("with offset with limit", func(t *testing.T) {
		testOffset := 8
		testLimit := 5
		err = Copy(fromFile.Name(), toFilePath, int64(testOffset), int64(testLimit))
		require.NoError(t, err)

		resultBytes := readFileBytes(t, toFilePath)
		require.Equal(t, testContentBytes[testOffset:], resultBytes)
	})
}

func readFileBytes(t *testing.T, toFilePath string) []byte {
	t.Helper()
	resultBytes, err := ioutil.ReadFile(toFilePath)
	require.NoError(t, err)
	return resultBytes
}

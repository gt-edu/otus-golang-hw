package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("validate filenames", func(t *testing.T) {
		dir := t.TempDir()
		defer cleanUpDir(dir)

		err := os.WriteFile(dir+"/Wrong=name", []byte(""), 0o777)
		require.NoError(t, err)

		envMap, err := ReadDir(dir)
		require.EqualError(t, err, ErrUnsupportedFilename.Error())
		require.Nil(t, envMap)
	})

	t.Run("Simple read", func(t *testing.T) {
		dir := t.TempDir()
		defer cleanUpDir(dir)

		testKey := "TEST"
		testValue := "test_value"
		err := os.WriteFile(dir+"/"+testKey, []byte(testValue+"\nSecond line"), 0o777)
		require.NoError(t, err)

		envMap, err := ReadDir(dir)
		require.NoError(t, err)
		require.NotNil(t, envMap)

		if val, ok := envMap[testKey]; ok {
			require.Equal(t, testValue, val.Value)
			require.Equal(t, false, val.NeedRemove)
		} else {
			require.Fail(t, "variable does not exists")
		}
	})

	t.Run("Replace special characters", func(t *testing.T) {
		zeroStr := string(byte(0))
		tests := []struct {
			filename      string
			firstLine     string
			expectedValue string
		}{
			{"TEST_TRIM1", "   value", "value"},
			{"TEST_TRIM2", "value   ", "value"},
			{"TEST_TRIM3", "\t\t\tvalue", "value"},
			{"TEST_TRIM4", "value\t\t", "value"},
			{"TEST_TRIM5", "   value\t\t   ", "value"},
			{"TEST_TRIM6", "\t\t   value ", "value"},
			{"TEST_REPLACE_ZERO", zeroStr + "first" + zeroStr + "second" + zeroStr, "\nfirst\nsecond\n"},
		}

		for _, tt := range tests {
			t.Run(tt.filename, func(t *testing.T) {
				dir := t.TempDir()
				defer cleanUpDir(dir)

				err := os.WriteFile(dir+"/"+tt.filename, []byte(tt.firstLine+"\nSecond line"), 0o777)
				require.NoError(t, err)

				envMap, err := ReadDir(dir)
				require.NoError(t, err)
				require.NotNil(t, envMap)

				if val, ok := envMap[tt.filename]; ok {
					require.Equal(t, tt.expectedValue, val.Value)
					require.Equal(t, false, val.NeedRemove)
				} else {
					require.Fail(t, "variable does not exists")
				}
			})
		}
	})

	t.Run("Empty files", func(t *testing.T) {
		tests := []struct {
			filename  string
			firstLine string
		}{
			{"TEST_EMPTY_FILE1", "   "},
			{"TEST_EMPTY_FILE3", "\t\t\t"},
			{"TEST_EMPTY_FILE4", "\t\t"},
			{"TEST_EMPTY_FILE5", "   \t\t   "},
			{"TEST_EMPTY_FILE6", "\t\t   "},
		}

		for _, tt := range tests {
			t.Run(tt.filename, func(t *testing.T) {
				dir := t.TempDir()
				defer cleanUpDir(dir)

				err := os.WriteFile(dir+"/"+tt.filename, []byte(tt.firstLine), 0o777)
				require.NoError(t, err)

				envMap, err := ReadDir(dir)
				require.NoError(t, err)
				require.NotNil(t, envMap)

				if val, ok := envMap[tt.filename]; ok {
					require.Equal(t, "", val.Value)
					require.Equal(t, true, val.NeedRemove)
				} else {
					require.Fail(t, "variable does not exists")
				}
			})
		}
	})
}

func cleanUpDir(dirName string) {
	if err := os.RemoveAll(dirName); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error during files removing: %v", err)
	}
}

package main

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("run with stdin and extra arguments", func(t *testing.T) {
		mockStdin, mockStdout, mockStderr := resetExecutorStreams()

		testStdinValue := "test stdin value"
		mockStdin.WriteString(testStdinValue)

		code := RunCmd([]string{"./testdata/test-cmd.sh", "0", "testarg"}, Environment{})
		require.Equal(t, 0, code)

		require.Empty(t, mockStderr.String())
		require.Equal(
			t, "VALUE: (); EXIT CODE: 0; All arguments: 0 testarg; Stdin content: "+testStdinValue,
			mockStdout.String(),
		)
	})

	t.Run("run with different exit code and stderr", func(t *testing.T) {
		_, mockStderr, mockStdout := resetExecutorStreams()

		code := RunCmd([]string{"./testdata/test-cmd.sh", "123"}, Environment{})
		require.Equal(t, 123, code)

		require.Equal(t, "VALUE: (); EXIT CODE: 123; All arguments: 123; Stdin content: ", mockStderr.String())
		require.Empty(t, mockStdout.String())
	})

	t.Run("run with different environment variables", func(t *testing.T) {
		_, mockStdout, mockStderr := resetExecutorStreams()

		_ = os.Setenv("HELLO", "SHOULD_REPLACE")
		_ = os.Setenv("FOO", "SHOULD_REPLACE")
		_ = os.Setenv("UNSET", "SHOULD_REMOVE")
		_ = os.Setenv("ADDED", "from original env")
		_ = os.Setenv("EMPTY", "SHOULD_BE_EMPTY")

		newEnv := Environment{
			"HELLO": EnvValue{Value: "\"hello\"", NeedRemove: false},
			"BAR":   EnvValue{Value: "bar", NeedRemove: false},
			"FOO":   EnvValue{Value: "   foo\nwith new line", NeedRemove: false},
			"UNSET": EnvValue{Value: "", NeedRemove: true},
			"EMPTY": EnvValue{Value: "", NeedRemove: true},
		}

		code := RunCmd([]string{"./testdata/echo.sh", "arg1=1", "arg2=2"}, newEnv)
		require.Equal(t, 0, code)

		require.Empty(t, mockStderr.String())
		require.Equal(t, `HELLO is ("hello")
BAR is (bar)
FOO is (   foo
with new line)
UNSET is ()
ADDED is (from original env)
EMPTY is ()
arguments are arg1=1 arg2=2
`, mockStdout.String())
	})

	// Different variables cases
}

func resetExecutorStreams() (*bytes.Buffer, *bytes.Buffer, *bytes.Buffer) {
	var mockStdin, mockStdout, mockStderr bytes.Buffer
	executorStreams = ExcecutorStreams{stderr: &mockStderr, stdout: &mockStdout, stdin: &mockStdin}
	return &mockStdin, &mockStdout, &mockStderr
}

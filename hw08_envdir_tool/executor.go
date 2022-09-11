package main

import (
	"errors"
	"io"
	"log"
	"os"
	exec "os/exec"
)

type ExcecutorStreams struct {
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
}

var executorStreams = ExcecutorStreams{stdin: os.Stdin, stdout: os.Stdout, stderr: os.Stderr}

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	cmdName := cmd[0]
	args := []string{}
	if len(cmd) > 1 {
		args = cmd[1:]
	}
	cmdExec := exec.Command(cmdName, args...)
	cmdExec.Stdin = executorStreams.stdin
	cmdExec.Stdout = executorStreams.stdout
	cmdExec.Stderr = executorStreams.stderr
	assembleEnvArray(env)

	if err := cmdExec.Start(); err != nil {
		log.Fatalf("cmd.Start: %v", err)
	}

	returnCode = 0
	if err := cmdExec.Wait(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			returnCode = exitErr.ExitCode()
		} else {
			log.Fatalf("cmd.Wait: %v", err)
		}
	}
	return
}

func assembleEnvArray(env Environment) {
	for key, element := range env {
		val := element.Value
		if element.NeedRemove {
			if _, exists := os.LookupEnv(key); exists {
				_ = os.Unsetenv(key)
			}
		} else {
			_ = os.Setenv(key, val)
		}
	}
}

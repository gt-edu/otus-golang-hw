package main

import (
	"io"
	"log"
	"os"
	"os/exec"
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
		if exitErr, isExitError := err.(*exec.ExitError); isExitError {
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
				if err := os.Unsetenv(key); err != nil {
					log.Printf("Could not unset variable: %v", err)
				}
			}
		} else {
			if err := os.Setenv(key, val); err != nil {
				log.Printf("Could not set variable: %v", err)
			}
		}
	}
}

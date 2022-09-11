package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

var ErrUnsupportedFilename = errors.New("unsupported file name")

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	env := make(Environment)
	for _, dirEntry := range dirEntries {
		if !dirEntry.IsDir() {
			if strings.Contains(dirEntry.Name(), "=") {
				return nil, ErrUnsupportedFilename
			}

			firstLine, err := getFirstLine(dir, dirEntry)
			if err != nil {
				return nil, err
			}

			firstLine = strings.TrimRight(firstLine, " \t")
			firstLine = strings.ReplaceAll(firstLine, string(byte(0)), "\n")

			needRemove := len(firstLine) == 0
			env[dirEntry.Name()] = EnvValue{
				Value:      firstLine,
				NeedRemove: needRemove,
			}
		}
	}

	return env, nil
}

func getFirstLine(dir string, dirEntry os.DirEntry) (string, error) {
	file, err := os.Open(filepath.Join(dir, dirEntry.Name()))
	if err != nil {
		return "", err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "error occurred during file closing: %v", err)
		}
	}(file)

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	firstLine := scanner.Text()
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return firstLine, nil
}

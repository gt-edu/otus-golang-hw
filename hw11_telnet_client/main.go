package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

func main() {
	address, t, err, flags := ParseTelnetClientFlags(os.Args)
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			PrintHelpAndExit(nil, flags)
		}
		PrintHelpAndExit(err, flags)
	}

	client := NewTelnetClient(address, t, os.Stdin, os.Stdout)
	err = client.Connect()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Connect error: %s\n", err.Error())
		os.Exit(1)
	}

	go func() {
		err := client.Receive()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
			return
		}
	}()
	err = client.Send()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
	}
}

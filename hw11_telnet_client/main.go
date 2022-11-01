package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
)

func main() {
	exitCode := StartTelnetClient()

	os.Exit(exitCode)
}

func StartTelnetClient() int {
	exitCode := 0
	address, t, flags, err := ParseTelnetClientFlags(os.Args)
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			PrintHelpAndExit(nil, flags)
		}
		PrintHelpAndExit(err, flags)
	}

	client := NewTelnetClient(address, t, os.Stdin, os.Stdout)
	err = client.Connect()
	if err != nil {
		PrintfToStderr("Connect error: %s\n", err.Error())
		os.Exit(1)
	}
	PrintfToStderr("Connected to %s.\n", address)
	defer func() {
		err := client.Close()
		if err != nil {
			PrintfToStderr("Error when final close: %s\n", err.Error())
		}
	}()

	go func() {
		err := client.Receive()
		if err != nil || client.ExitedWithEOF() {
			switch {
			case client.ExitedWithEOF():
				PrintfToStderr("...Connection was closed by peer\n")
			case errors.Is(err, net.ErrClosed):
				PrintfToStderr("...Connection was closed on interruption\n")
			default:
				PrintfToStderr("Error when receive: %s\n", err.Error())
			}
			err := client.Close()
			if err != nil {
				PrintfToStderr("Error when close connection after receive: %s\n", err.Error())
			}

			return
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	go func() {
		defer stop()
		<-ctx.Done()
		err := client.Close()
		if err != nil {
			PrintfToStderr("Error when closing client on interrupt: %s\n", err.Error())
			return
		}
	}()

	err = client.Send()
	if err != nil || client.ExitedWithEOF() {
		switch {
		case client.ExitedWithEOF():
			PrintfToStderr("...EOF\n")
		default:
			PrintfToStderr("Error when send: %s\n", err.Error())
			exitCode = 1
		}
	}
	return exitCode
}

func PrintfToStderr(msg string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, msg, args...)
}

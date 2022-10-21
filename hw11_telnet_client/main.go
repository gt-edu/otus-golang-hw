package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
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
		PrintfToStderr("Connect error: %s\n", err.Error())
		os.Exit(1)
	}
	PrintfToStderr("Connected to %s.\n", address)

	go func() {
		err := client.Receive()
		if err != nil {
			switch {
			case errors.Is(err, io.EOF):
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
	defer stop()
	go func() {
		<-ctx.Done()
		err := client.Close()
		if err != nil {
			PrintfToStderr("Error when closing client on interrupt: %s\n", err.Error())
			return
		}
	}()

	err = client.Send()
	if err != nil {
		switch {
		case errors.Is(err, io.EOF):
			PrintfToStderr("...EOF\n")
		default:
			PrintfToStderr("Error when send: %s\n", err.Error())
			os.Exit(1)
		}
	}
}

func PrintfToStderr(msg string, args ...interface{}) {
	if len(args) > 0 {
		_, _ = fmt.Fprintf(os.Stderr, msg, args)
	} else {
		_, _ = fmt.Fprintf(os.Stderr, msg)
	}

}

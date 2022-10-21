package main

import (
	"bytes"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

var (
	ErrInvalidPort                = errors.New("invalid port value")
	ErrHostAndPortMustBeSpecified = errors.New("host and port must be specified")
	ErrParseError                 = errors.New("parse error")
)

func ParseTelnetClientFlags(args []string) (string, time.Duration, error, *flag.FlagSet) {

	flags := flag.NewFlagSet(args[0], flag.ContinueOnError)
	var buf bytes.Buffer
	flags.SetOutput(&buf)

	var timeout time.Duration
	flags.DurationVar(&timeout, "timeout", time.Duration(10)*time.Second, "Specifies connection timeout")

	err := flags.Parse(args[1:])
	if err != nil {
		if err == flag.ErrHelp {
			return "", 0, err, flags
		} else {
			return "", 0, errors.Wrap(ErrParseError, err.Error()), flags
		}
	}

	tail := flags.Args()
	if len(tail) < 2 {
		return "", timeout, ErrHostAndPortMustBeSpecified, flags
	}

	port := tail[1]
	_, err = strconv.Atoi(port)
	if err != nil {
		return "", 0, errors.Wrap(ErrInvalidPort, err.Error()), flags
	}

	return net.JoinHostPort(tail[0], port), timeout, nil, flags
}

func PrintHelpAndExit(err error, flags *flag.FlagSet) {
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
	}
	_, _ = fmt.Fprintf(os.Stderr, "Usage of %s: [parameters] host port\n", os.Args[0])

	flags.PrintDefaults()
	os.Exit(1)
}

package main

import (
	"flag"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestParseTelnetClientFlags(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		wantAddress string
		wantSeconds int
		wantErr     error
	}{
		{"timeout 10s", []string{"", "--timeout=9s", "localhost", "8080"}, "localhost:8080", 9, nil},
		{"no timeout", []string{"", "localhost", "8080"}, "localhost:8080", 10, nil},
		{"help", []string{"", "-h"}, "", 0, flag.ErrHelp},
		{"empty", []string{""}, "", 10, ErrHostAndPortMustBeSpecified},
		{"timeout incorrect", []string{"", "--timeout=zs", "localhost", "8080"}, "", 0, ErrParseError},
		{"port incorrect", []string{"", "localhost", "zz"}, "", 0, ErrInvalidPort},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			address, timeout, _, err := ParseTelnetClientFlags(tt.args)
			require.ErrorIs(t, err, tt.wantErr)
			require.Equal(t, tt.wantAddress, address)
			require.Equal(t, time.Duration(tt.wantSeconds)*time.Second, timeout)
		})
	}
}

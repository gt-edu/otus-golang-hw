package main

import (
	"flag"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestParseTelnetClientFlags(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		wantAddress string
		wantSeconds int
		wantErr     error
	}{
		{"timeout 10s", []string{"", "--timeout=10s", "localhost", "8080"}, "localhost:8080", 10, nil},
		{"no timeout", []string{"", "localhost", "8080"}, "localhost:8080", 0, nil},
		{"help", []string{"", "-h"}, "", 0, flag.ErrHelp},
		{"empty", []string{""}, "", 0, ErrHostAndPortMustBeSpecified},
		{"timeout incorrect", []string{"", "--timeout=zs", "localhost", "8080"}, "", 0, ErrParseError},
		{"port incorrect", []string{"", "localhost", "zz"}, "", 0, ErrInvalidPort},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			address, timeout, err, _ := ParseTelnetClientFlags(tt.args)
			require.ErrorIs(t, err, tt.wantErr)
			require.Equal(t, tt.wantAddress, address)
			require.Equal(t, time.Duration(tt.wantSeconds)*time.Second, timeout)
		})
	}
}

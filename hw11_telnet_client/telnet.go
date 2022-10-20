package main

import (
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	// Place your code here.
	return &TelnetClientImpl{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

type TelnetClientImpl struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer

	// Нужно ли здесь указатель
	conn net.Conn
}

func (t *TelnetClientImpl) Connect() error {
	var err error
	t.conn, err = net.Dial("tcp", t.address)
	if err != nil {
		return err
	}

	return nil
}

func (t *TelnetClientImpl) Close() error {
	if t.conn != nil {
		err := t.conn.Close()
		if err != nil {
			return err
		}
		t.conn = nil
	}

	return nil
}

func (t *TelnetClientImpl) Send() error {
	p := make([]byte, 10000)
	for {
		n, err := t.in.Read(p)
		if err != nil {
			return err
		}
		_, err = t.conn.Write(p[:n])
		if err != nil {
			return err
		}
	}
}

func (t *TelnetClientImpl) Receive() error {
	p := make([]byte, 10000)
	for {
		n, err := t.conn.Read(p)
		if err != nil {
			return err
		}

		_, err = t.out.Write(p[:n])
		if err != nil {
			return err
		}
	}
}

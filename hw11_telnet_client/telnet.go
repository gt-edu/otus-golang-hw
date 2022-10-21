package main

import (
	"fmt"
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

	return &TelnetClientImpl{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,

		doneCh: make(chan struct{}),
	}
}

type TelnetClientImpl struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer

	// Нужно ли здесь указатель
	conn net.Conn

	doneCh chan struct{}
	closed bool
}

func (t *TelnetClientImpl) Connect() error {
	var err error
	t.conn, err = net.DialTimeout("tcp", t.address, t.timeout)
	if err != nil {
		return err
	}

	return nil
}

func (t *TelnetClientImpl) Close() error {
	if t.closed {
		return nil
	}
	t.closed = true

	close(t.doneCh)

	err := t.in.Close()
	if err != nil {
		return err
	}

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
		select {
		case <-t.doneCh:
			return nil
		default:
			n, err := t.in.Read(p)
			if err != nil {
				return err
			}

			if t.conn != nil {
				_, err = t.conn.Write(p[:n])
				if err != nil {
					return err
				}
			}
		}
	}
}

func (t *TelnetClientImpl) Receive() error {
	p := make([]byte, 10000)
	for {
		select {
		case <-t.doneCh:
			fmt.Println("Gracefully exit when receiving")
			return nil
		default:
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
}

package dnet

import (
	"errors"
	"io"
	"net"
	"strings"
	"sync/atomic"
	"time"
)

type TCPAcceptor struct {
	address  string
	listener net.Listener
	started  int32
}

// NewTCPAcceptor returns a new instance of TCPAcceptor
func NewTCPAcceptor(address string) *TCPAcceptor {
	return &TCPAcceptor{address: address}
}

// ServeTCP listen and serve tcp address with AcceptorHandler
func ServeTCP(address string, handler AcceptorHandler) error {
	return NewTCPAcceptor(address).Serve(handler)
}

// ServeTCPFunc listen and serve tcp address with AcceptorHandlerFunc
func ServeTCPFunc(address string, handler AcceptorHandlerFunc) error {
	return NewTCPAcceptor(address).ServeFunc(handler)
}

// Serve listens and serve in the specified addr
func (this *TCPAcceptor) Serve(handler AcceptorHandler) error {
	if handler == nil {
		return errors.New("dnet:Serve handler is nil. ")
	}

	if !atomic.CompareAndSwapInt32(&this.started, 0, 1) {
		return errors.New("dnet:Serve acceptor is already started. ")
	}

	listener, err := net.Listen("tcp", this.address)
	if err != nil {
		return err
	}
	this.listener = listener
	defer this.Stop()

	var tempDelay time.Duration
	for {
		conn, err := this.listener.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				time.Sleep(tempDelay)
				continue
			}
			if strings.Contains(err.Error(), "use of closed network connection") {
				return io.EOF
			}
			return err
		}

		go handler.OnConnection(conn)
	}

}

// ServeFunc listens and serve in the specified addr
func (this *TCPAcceptor) ServeFunc(handler AcceptorHandlerFunc) error {
	return this.Serve(handler)
}

// Addr returns the addr the acceptor will listen on
func (this *TCPAcceptor) Addr() net.Addr {
	return this.listener.Addr()
}

// Stop stops the acceptor
func (this *TCPAcceptor) Stop() {
	if atomic.CompareAndSwapInt32(&this.started, 1, 0) {
		_ = this.listener.Close()
	}
}

// DialTCP
func DialTCP(address string, timeout time.Duration) (net.Conn, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return nil, err
	}

	dialer := &net.Dialer{Timeout: timeout}
	return dialer.Dial(tcpAddr.Network(), address)
}

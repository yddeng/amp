package dnet

import (
	"errors"
	"github.com/gorilla/websocket"
	"log"
	"net"
	"net/http"
	"net/url"
	"sync/atomic"
	"time"
)

type WSAcceptor struct {
	address  string
	handler  *wsHandler
	listener net.Listener
	started  int32
}

// NewWSAcceptor returns a new instance of WSAcceptor
func NewWSAcceptor(address string) *WSAcceptor {
	return &WSAcceptor{
		address: address,
		handler: &wsHandler{
			upgrader: &websocket.Upgrader{
				CheckOrigin: func(r *http.Request) bool {
					// allow all connections by default
					return true
				},
			},
		},
	}
}

// ServeWS listen and serve ws address with AcceptorHandler
func ServeWS(address string, handler AcceptorHandler) error {
	return NewWSAcceptor(address).Serve(handler)
}

// ServeWS listen and serve ws address with AcceptorHandlerFunc
func ServeWSFunc(address string, handler AcceptorHandlerFunc) error {
	return NewWSAcceptor(address).Serve(handler)
}

type wsHandler struct {
	upgrader *websocket.Upgrader
	handler  AcceptorHandler
}

func (h *wsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("dnet:ServeHTTP WSSession Upgrade failed, %s\n", err.Error())
		return
	}
	h.handler.OnConnection(NewWSConn(c))
}

// Serve listens and serve in the specified addr
func (this *WSAcceptor) Serve(handler AcceptorHandler) error {
	if handler == nil {
		return errors.New("dnet:Serve handler is nil. ")
	}
	this.handler.handler = handler

	if !atomic.CompareAndSwapInt32(&this.started, 0, 1) {
		return errors.New("dnet:Serve acceptor is already started. ")
	}

	listener, err := net.Listen("tcp", this.address)
	if err != nil {
		return errors.New("dnet:Serve net.Listen failed, " + err.Error())
	}
	this.listener = listener
	defer this.Stop()

	if err = http.Serve(this.listener, this.handler); err != nil {
		log.Printf("dnet:Serve failed, %s\n", err.Error())
	}

	return nil
}

// ServeFunc listens and serve in the specified addr
func (this *WSAcceptor) ServeFunc(handler AcceptorHandlerFunc) error {
	return this.Serve(handler)
}

// Addr returns the addr the acceptor will listen on
func (this *WSAcceptor) Addr() net.Addr {
	return this.listener.Addr()
}

// Stop stops the acceptor
func (this *WSAcceptor) Stop() {
	if !atomic.CompareAndSwapInt32(&this.started, 1, 0) {
		_ = this.listener.Close()
	}
}

func DialWS(host string, timeout time.Duration) (net.Conn, error) {
	u := url.URL{Scheme: "ws", Host: host}
	websocket.DefaultDialer.HandshakeTimeout = timeout
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, err
	}
	return NewWSConn(conn), nil
}

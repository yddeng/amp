package dnet

import (
	"github.com/gorilla/websocket"
	"io"
	"net"
	"time"
)

// WSConn is an adapter to t.Conn, which implements all t.Conn
// interface base on *websocket.Conn
type WSConn struct {
	conn   *websocket.Conn
	typ    int // message type
	reader io.Reader
}

// NewWSConn return an initialized *WSConn
func NewWSConn(conn *websocket.Conn) *WSConn {
	return &WSConn{conn: conn}
}

// Read reads data from the connection.
// Read can be made to time out and return an Error with Timeout() == true
// after a fixed time limit; see SetDeadline and SetReadDeadline.
func (c *WSConn) Read(b []byte) (int, error) {
reRead:
	if c.reader == nil {
		t, r, err := c.conn.NextReader()
		if err != nil {
			return 0, err
		}
		c.typ = t
		c.reader = r
	}

	off, lenb := 0, len(b)
	if lenb == 0 {
		return 0, nil
	}
	for off < lenb {
		n, err := c.reader.Read(b[off:])
		if err != nil {
			if err == io.EOF {
				c.reader = nil
				err = nil
				if off == 0 {
					goto reRead
				}
			}
			return off + n, err
		}
		off += n
	}
	return off, nil
}

// Write writes data to the connection.
// Write can be made to time out and return an Error with Timeout() == true
// after a fixed time limit; see SetDeadline and SetWriteDeadline.
func (c *WSConn) Write(b []byte) (int, error) {
	err := c.conn.WriteMessage(websocket.BinaryMessage, b)
	if err != nil {
		return 0, err
	}

	return len(b), nil
}

// Close closes the connection.
// Any blocked Read or Write operations will be unblocked and return errors.
func (c *WSConn) Close() error {
	return c.conn.Close()
}

// LocalAddr returns the local network address.
func (c *WSConn) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

// RemoteAddr returns the remote network address.
func (c *WSConn) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

// SetDeadline sets the read and write deadlines associated
// with the connection. It is equivalent to calling both
// SetReadDeadline and SetWriteDeadline.
//
// A deadline is an absolute time after which I/O operations
// fail with a timeout (see type Error) instead of
// blocking. The deadline applies to all future and pending
// I/O, not just the immediately following call to Read or
// Write. After a deadline has been exceeded, the connection
// can be refreshed by setting a deadline in the future.
//
// An idle timeout can be implemented by repeatedly extending
// the deadline after successful Read or Write calls.
//
// A zero value for t means I/O operations will not time out.
func (c *WSConn) SetDeadline(t time.Time) error {
	if err := c.SetReadDeadline(t); err != nil {
		return err
	}

	return c.SetWriteDeadline(t)
}

// SetReadDeadline sets the deadline for future Read calls
// and any currently-blocked Read call.
// A zero value for t means Read will not time out.
func (c *WSConn) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

// SetWriteDeadline sets the deadline for future Write calls
// and any currently-blocked Write call.
// Even if write times out, it may return n > 0, indicating that
// some of the data was successfully written.
// A zero value for t means Write will not time out.
func (c *WSConn) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}

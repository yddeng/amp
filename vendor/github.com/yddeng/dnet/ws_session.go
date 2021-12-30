package dnet

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"reflect"
)

type DefWsCodec struct{}

//解码
func (_ DefWsCodec) Decode(reader io.Reader) (interface{}, error) {
	buff := new(bytes.Buffer)
	_, err := buff.ReadFrom(reader)
	if err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}

//编码
func (_ DefWsCodec) Encode(o interface{}) ([]byte, error) {
	data, ok := o.([]byte)
	if !ok {
		return nil, fmt.Errorf("dnet:defWSCodec encode interface{} is %s, need type []byte", reflect.TypeOf(o))
	}
	return data, nil
}

type WSSession struct {
	*session
}

// NewWSSession return an initialized *WSSession
func NewWSSession(conn net.Conn, options ...Option) *WSSession {
	op := loadOptions(options...)
	if op.MsgCallback == nil {
		// need message callback
		panic(ErrNilMsgCallBack)
	}
	// init default codec
	if op.Codec == nil {
		op.Codec = DefWsCodec{}
	}

	return &WSSession{
		session: newSession(conn, op),
	}
}

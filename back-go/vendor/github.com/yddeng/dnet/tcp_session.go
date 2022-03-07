package dnet

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"reflect"
)

// default编解码器
// 消息 -- 格式: 消息头(消息len), 消息体

const (
	lenSize  = 2       // 消息长度（消息体的长度）
	headSize = lenSize // 消息头长度
	buffSize = 65535   // 缓存容量(与lenSize有关，2字节最大65535）
)

type DefTCPCodec struct{}

//解码
func (_ DefTCPCodec) Decode(reader io.Reader) (interface{}, error) {
	hdr := make([]byte, headSize)
	_, err := io.ReadFull(reader, hdr)
	if err != nil {
		return nil, err
	}

	length := binary.BigEndian.Uint16(hdr)
	buff := make([]byte, length)

	if _, err := io.ReadFull(reader, buff); err != nil {
		return nil, err
	}
	return buff, nil
}

//编码
func (_ DefTCPCodec) Encode(o interface{}) ([]byte, error) {

	data, ok := o.([]byte)
	if !ok {
		return nil, fmt.Errorf("dnet:Encode interface{} is %s, need type []byte", reflect.TypeOf(o))
	}

	length := len(data)
	if length > buffSize {
		return nil, fmt.Errorf("dnet:Encode dataLen is too large,len: %d", length)
	}

	buff := new(bytes.Buffer)
	binary.Write(buff, binary.BigEndian, uint16(length))
	buff.Write(data)
	return buff.Bytes(), nil
}

// TCPSession
type TCPSession struct {
	*session
}

// NewTCPSession return an initialized *TCPSession
func NewTCPSession(conn net.Conn, options ...Option) *TCPSession {
	op := loadOptions(options...)
	if op.MsgCallback == nil {
		// need message callback
		panic(ErrNilMsgCallBack)
	}
	// init default codec
	if op.Codec == nil {
		op.Codec = DefTCPCodec{}
	}

	return &TCPSession{
		session: newSession(conn, op),
	}
}

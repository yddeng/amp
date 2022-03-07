package protocol

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/yddeng/dnet/drpc"
	"io"
	"reflect"
)

const (
	flagSize = 1
	seqSize  = 8
	cmdSize  = 2
	bodySize = 4
	HeadSize = flagSize + seqSize + cmdSize + bodySize
	BuffSize = 1024 * 256
)

type flag byte

const (
	message flag = 0x80
	rpcReq  flag = 0x40
	rpcResp flag = 0x20
)

func (this flag) getType() flag {
	if this&message > 0 {
		return message
	} else if this&rpcReq > 0 {
		return rpcReq
	} else if this&rpcResp > 0 {
		return rpcResp
	}
	return message
}

type Codec struct{}

func readHeader(buff []byte) (byte, uint64, uint16, uint32) {
	var b byte
	var seq uint64
	var length uint32
	var cmd uint16

	b = buff[0]
	buffer := bytes.NewBuffer(buff[1:])
	binary.Read(buffer, binary.BigEndian, &seq)
	binary.Read(buffer, binary.BigEndian, &cmd)
	binary.Read(buffer, binary.BigEndian, &length)

	return b, seq, cmd, length
}

//解码
func (*Codec) Decode(reader io.Reader) (interface{}, error) {
	hdr := make([]byte, HeadSize)
	if _, err := io.ReadFull(reader, hdr); err != nil {
		return nil, err
	}

	b, seq, cmd, length := readHeader(hdr)
	if length < 0 || length >= BuffSize {
		return nil, fmt.Errorf("Message too large. ")
	}

	buff := make([]byte, length)
	if _, err := io.ReadFull(reader, buff); err != nil {
		return nil, err
	}

	tt := flag(b).getType()
	switch tt {
	case message:
		m, err := Unmarshal("msg", cmd, buff)
		if err != nil {
			return nil, err
		}
		return &Message{
			data: m.(proto.Message),
			cmd:  cmd,
		}, nil
	case rpcReq:
		m, err := Unmarshal("req", cmd, buff)
		if err != nil {
			return nil, err
		}
		return &drpc.Request{
			Seq:    seq,
			Method: proto.MessageName(m.(proto.Message)),
			Data:   m,
		}, nil
	case rpcResp:
		m, err := Unmarshal("resp", cmd, buff)
		if err != nil {
			return nil, err
		}
		return &drpc.Response{
			Seq:  seq,
			Data: m,
		}, nil
	default:
		return nil, errors.New("invalid")
	}
}

//编码
func (*Codec) Encode(o interface{}) ([]byte, error) {
	var b flag
	var seq uint64
	var cmd uint16
	var data []byte
	var err error

	switch o.(type) {
	case *Message:
		b = message
		msg := o.(*Message)
		cmd, data, err = Marshal("msg", msg.GetData())
		if err != nil {
			return nil, err
		}
	case *drpc.Request:
		msg := o.(*drpc.Request)
		b = rpcReq
		seq = msg.Seq
		cmd, data, err = Marshal("req", msg.Data)
		if err != nil {
			return nil, err
		}
	case *drpc.Response:
		msg := o.(*drpc.Response)
		b = rpcResp
		seq = msg.Seq
		cmd, data, err = Marshal("resp", msg.Data)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("invalid type:%s. ", reflect.TypeOf(o).String())
	}

	length := uint32(len(data))
	buffer := new(bytes.Buffer)
	buffer.WriteByte(byte(b))
	binary.Write(buffer, binary.BigEndian, seq)
	binary.Write(buffer, binary.BigEndian, cmd)
	binary.Write(buffer, binary.BigEndian, length)
	buffer.Write(data)

	return buffer.Bytes(), nil
}

type Message struct {
	data interface{}
	cmd  uint16
}

func NewMessage(data interface{}) *Message {
	msg := &Message{
		data: data,
	}
	return msg
}

func (this *Message) GetData() interface{} {
	return this.data
}

func (this *Message) GetCmd() uint16 {
	return this.cmd
}

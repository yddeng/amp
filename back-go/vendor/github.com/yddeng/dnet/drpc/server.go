package drpc

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
)

type Request struct {
	Seq    uint64 // the number of request
	Method string // The name of the service and method to call.
	Data   interface{}
}

type Response struct {
	Seq   uint64 // the number of request
	Data  interface{}
	Error string
}

// RPCChannel
type RPCChannel interface {
	SendRequest(req *Request) error    // send rpc request
	SendResponse(resp *Response) error // send rpc response
}

// Server represents an RPC Server.
type Server struct {
	methods map[string]MethodHandler
	mtx     sync.RWMutex
}

type MethodHandler func(replier *Replier, req interface{})

// Register Register the method on the server whit name.
// call method by name.
func (server *Server) Register(name string, h MethodHandler) {
	if name == "" {
		panic("drpc: Register name == ''")
	}
	if nil == h {
		panic("drpc: Register h == nil")
	}

	server.mtx.Lock()
	defer server.mtx.Unlock()
	_, ok := server.methods[name]
	if ok {
		panic(fmt.Sprintf("drpc:Register duplicate method:%s", name))
	}
	server.methods[name] = h
}

// OnRPCRequest
func (server *Server) OnRPCRequest(channel RPCChannel, req *Request) error {
	if channel == nil || req == nil {
		return fmt.Errorf("drpc:OnRPCRequest invalid argument")
	}

	server.mtx.RLock()
	method, ok := server.methods[req.Method]
	server.mtx.RUnlock()
	if !ok {
		return fmt.Errorf("drpc:OnRPCRequest invalid method %s", req.Method)
	}

	replier := &Replier{Channel: channel, resp: &Response{Seq: req.Seq}}
	return server.callMethod(method, replier, req.Data)
}

func (server *Server) callMethod(method MethodHandler, replier *Replier, arg interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, 1024)
			l := runtime.Stack(buf, false)
			err = fmt.Errorf(fmt.Sprintf("%v: %s", r, buf[:l]))
		}
	}()

	method(replier, arg)
	return
}

// Replier
type Replier struct {
	Channel RPCChannel
	fired   int32
	resp    *Response
}

func (r *Replier) Reply(ret interface{}, err error) error {
	if !atomic.CompareAndSwapInt32(&r.fired, 0, 1) {
		return fmt.Errorf("drpc:Reply repeated reply %d ", atomic.LoadInt32(&r.fired))
	}

	if err != nil {
		r.resp.Error = err.Error()
	} else if ret != nil {
		r.resp.Data = ret
	} else {
		return fmt.Errorf("drpc:Reply argments failed, none")
	}
	return r.reply(r.resp)
}

func (r *Replier) reply(resp *Response) error {
	return r.Channel.SendResponse(resp)
}

// NewServer returns a new Server.
func NewServer() *Server {
	return &Server{
		methods: map[string]MethodHandler{},
	}
}

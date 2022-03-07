package drpc

import (
	"errors"
	"fmt"
	"github.com/yddeng/timer"
	"sync"
	"sync/atomic"
	"time"
)

const DefaultRPCTimeout = 8 * time.Second

var ErrRPCTimeout = fmt.Errorf("drpc: rpc timeout. ")

// Call represents an active RPC.
type Call struct {
	reqNo    uint64
	callback func(interface{}, error)
	timer    timer.Timer
}

// Client represents an RPC Client.
// There may be multiple outstanding Calls associated
// with a single Client, and a Client may be used by
// multiple goroutines simultaneously.
type Client struct {
	reqNo    uint64         // serial number
	timerMgr timer.TimerMgr // timer
	pending  sync.Map       //map[uint64]*Call
}

// Call invokes the function synchronous, waits for it to complete, and returns its result and error status.
func (client *Client) Call(channel RPCChannel, method string, data interface{}, timeout time.Duration) (result interface{}, err error) {
	waitC := make(chan struct{})
	f := func(ret_ interface{}, err_ error) {
		result = ret_
		err = err_
		close(waitC)
	}
	if err := client.Go(channel, method, data, timeout, f); err != nil {
		return nil, err
	}
	<-waitC
	return
}

// Go invokes the function asynchronously.
func (client *Client) Go(channel RPCChannel, method string, data interface{}, timeout time.Duration, callback func(interface{}, error)) error {
	if callback == nil {
		return fmt.Errorf("drpc: Go callback == nil")
	}

	seq := atomic.AddUint64(&client.reqNo, 1)
	c := &Call{reqNo: seq, callback: callback}
	req := &Request{Seq: seq, Method: method, Data: data}
	client.pending.Store(seq, c)

	c.timer = client.timerMgr.OnceTimer(timeout, func() {
		if v, ok := client.pending.LoadAndDelete(seq); ok {
			v.(*Call).callback(nil, ErrRPCTimeout)
		}
	})

	if err := channel.SendRequest(req); err != nil {
		if v, ok := client.pending.LoadAndDelete(seq); ok {
			if v.(*Call).timer != nil {
				v.(*Call).timer.Stop()
			}
		}
		return err
	}

	//client.lock.Lock()
	//seq := client.reqNo
	//client.reqNo++
	//c := &Call{reqNo: seq, callback: callback}
	//client.pending[seq] = c
	//client.lock.Unlock()
	//
	//c.timer = client.timerMgr.OnceTimer(timeout, func() {
	//	client.lock.Lock()
	//	if call, ok := client.pending[seq]; ok {
	//		delete(client.pending, seq)
	//		call.callback(nil, ErrRPCTimeout)
	//		call.timer = nil
	//	}
	//	client.lock.Unlock()
	//})
	//
	//req := &Request{Seq: seq, Method: method, Data: data}
	//
	//// 避免channel 中直接调用 OnRPCResponse, 导致死锁
	//if err := channel.SendRequest(req); err != nil {
	//	client.lock.Lock()
	//	if call, ok := client.pending[seq]; ok {
	//		delete(client.pending, seq)
	//		if call.timer != nil {
	//			call.timer.Stop()
	//		}
	//	}
	//	client.lock.Unlock()
	//	return err
	//}

	return nil
}

// OnRPCResponse
func (client *Client) OnRPCResponse(resp *Response) error {
	v, ok := client.pending.LoadAndDelete(resp.Seq)
	if !ok {
		return fmt.Errorf("drpc: OnRPCResponse reqNo:%d is not found", resp.Seq)
	}

	call := v.(*Call)
	if resp.Error != "" {
		call.callback(nil, errors.New(resp.Error))
	} else {
		call.callback(resp.Data, nil)
	}

	if call.timer != nil {
		call.timer.Stop()
	}
	return nil

}

// NewClient returns a new Client to handle requests to the
// set of services at the other end of the connection.
// It adds a timer manager to
func NewClient() *Client {
	return &Client{
		timerMgr: timer.NewTimeWheelMgr(time.Millisecond*50, 200),
	}
}

// NewClientWithTimerMgr is like NewClient but uses the specified timerMgr.
func NewClientWithTimerMgr(timerMgr timer.TimerMgr) *Client {
	return &Client{
		timerMgr: timerMgr,
	}
}

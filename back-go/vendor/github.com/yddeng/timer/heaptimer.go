// MIT License
//
// Copyright (c) 2021 ydd
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package timer

import (
	"github.com/yddeng/utils/heap"
	"sync"
	"sync/atomic"
	"time"
)

// HeapTimer
type HeapTimer struct {
	key int64
	rt  *runtimeTimer
	mgr *HeapTimerMgr
}

// newHeapTimer creates an initialized heapTimer.
func newHeapTimer(d time.Duration, repeated bool, f func()) *HeapTimer {
	return &HeapTimer{
		key: 0,
		rt: &runtimeTimer{
			expire:   when(d),
			fn:       f,
			repeated: repeated,
			period:   d,
		},
	}
}

// Less is used to compare expiration with other Timer.
func (t *HeapTimer) Less(elem interface{}) bool {
	return t.rt.expire.Before(elem.(*HeapTimer).rt.expire)
}

// Stop prevents the Timer from firing.
// It returns true if the call stops the timer, false if the timer has already
// been stopped.
func (t *HeapTimer) Stop() bool {
	if t.mgr == nil {
		panic("timer: Stop called on uninitialized HeapTimer")
	}
	if atomic.CompareAndSwapInt32(&t.rt.stopped, 0, 1) {
		t.mgr.removeTimer(t)
		return true
	}
	return false
}

// Reset changes the timer to expire after duration d.
// It returns true if the timer had been active or had executed,
// false if the timer been stopped.
func (t *HeapTimer) Reset(d time.Duration) bool {
	if t.mgr == nil {
		panic("timer: Reset called on uninitialized HeapTimer")
	}

	if atomic.LoadInt32(&t.rt.stopped) == 1 {
		return false
	}

	t.mgr.removeTimer(t)
	t.rt.expire = when(d)
	t.rt.period = d
	t.mgr.addTimer(t)
	return true
}

// do calls f in its own goroutine. it return when it is stopped.
// addTimer when t.repeated is true.
func (t *HeapTimer) do() {
	if atomic.LoadInt32(&t.rt.stopped) == 1 {
		return
	}

	goFunc(t.rt.fn)
	//repeat
	if t.rt.repeated {
		if atomic.LoadInt32(&t.rt.stopped) == 1 {
			return
		}
		t.rt.expire = when(t.rt.period)
		t.mgr.addTimer(t)
	}

}

type HeapTimerMgr struct {
	minHeap     *heap.Heap
	accumulator int64 // accumulator for Timer.key
	timers      map[int64]*HeapTimer
	signalChan  chan struct{} // signal channel for add timer
	started     int32
	mtx         sync.Mutex
}

// NewHeapTimerMgr creates an initialized HeapTimerMgr.
func NewHeapTimerMgr() *HeapTimerMgr {
	return &HeapTimerMgr{
		minHeap:    heap.New(),
		timers:     map[int64]*HeapTimer{},
		signalChan: make(chan struct{}, 1),
	}
}

// addTimer add timer to manager
func (mgr *HeapTimerMgr) addTimer(t *HeapTimer) {
	if t.key == 0 {
		key := atomic.AddInt64(&mgr.accumulator, 1)
		t.key = key
		t.mgr = mgr
	}

	mgr.mtx.Lock()
	defer mgr.mtx.Unlock()
	mgr.timers[t.key] = t
	mgr.minHeap.Push(t)
	sendSignal(mgr.signalChan)
}

// removeTimer remove timer from manager
func (mgr *HeapTimerMgr) removeTimer(t *HeapTimer) {
	mgr.mtx.Lock()
	defer mgr.mtx.Unlock()
	if _, ok := mgr.timers[t.key]; ok {
		delete(mgr.timers, t.key)
		mgr.minHeap.Remove(t)
	}
}

// run
func (mgr *HeapTimerMgr) run() {
	var tt *time.Timer
	var min heap.Element
	for {
		// 每添加一个 timer 都会触发循环，重新计算最小值
		<-mgr.signalChan
		now := time.Now()
		for {
			mgr.mtx.Lock()
			min = mgr.minHeap.Top()
			mgr.mtx.Unlock()
			if nil != min && !now.Before(min.(*HeapTimer).rt.expire) {
				t := min.(*HeapTimer)
				mgr.removeTimer(t)
				t.do()
			} else {
				break
			}
		}

		if min != nil {
			sleepTime := min.(*HeapTimer).rt.expire.Sub(now)
			if nil == tt {
				tt = time.AfterFunc(sleepTime, func() {
					sendSignal(mgr.signalChan)
				})
			} else {
				tt.Reset(sleepTime)
			}
		}
	}

}

// registerTimer add timer to manager and run the manager first time
func (mgr *HeapTimerMgr) registerTimer(d time.Duration, repeated bool, f func()) *HeapTimer {
	if atomic.CompareAndSwapInt32(&mgr.started, 0, 1) {
		go mgr.run()
	}
	t := newHeapTimer(d, repeated, f)
	mgr.addTimer(t)
	return t
}

// OnceTimer waits for the duration to elapse and then calls f
// in its own goroutine. It returns a Timer. It's done once
func (mgr *HeapTimerMgr) OnceTimer(d time.Duration, f func()) Timer {
	return mgr.registerTimer(d, false, f)
}

// RepeatTimer waits for the duration to elapse and then calls f
// in its own goroutine. It returns a Timer. It can be used to
// cancel the call using its Stop method.
func (mgr *HeapTimerMgr) RepeatTimer(d time.Duration, f func()) Timer {
	return mgr.registerTimer(d, true, f)
}

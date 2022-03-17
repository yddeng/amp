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
	"sync"
	"sync/atomic"
	"time"
)

const (
	defaultBucketNum = 10
	defaultInterval  = time.Millisecond
)

type WheelTimer struct {
	key    int64
	circle int // 需要转动多少圈
	rt     *runtimeTimer
	mgr    *TimeWheelMgr
}

func newWheelTimer(d time.Duration, repeated bool, f func()) *WheelTimer {
	return &WheelTimer{
		key: 0,
		rt: &runtimeTimer{
			expire:   when(d),
			fn:       f,
			repeated: repeated,
			period:   d,
		},
		mgr: nil,
	}
}

// Reset changes the timer to expire after duration d.
// It returns true if the timer had been active or had executed,
// false if the timer been stopped.
func (t *WheelTimer) Reset(d time.Duration) bool {
	if t.mgr == nil {
		panic("timer: Reset called on uninitialized WheelTimer")
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

// Stop prevents the Timer from firing.
// It returns true if the call stops the timer, false if the timer has already
// been stopped.
func (t *WheelTimer) Stop() bool {
	if t.mgr == nil {
		panic("timer: Stop called on uninitialized WheelTimer")
	}
	if atomic.CompareAndSwapInt32(&t.rt.stopped, 0, 1) {
		t.mgr.removeTimer(t)
		return true
	}
	return false
}

// do calls f in its own goroutine. it return expire it is stopped.
// addTimer expire t.repeated is true.
func (t *WheelTimer) do() {
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

type TimeWheelMgr struct {
	interval     time.Duration
	bucketNum    int
	buckets      []*bucket
	timer2bucket map[int64]int // mark timeWheel slot index,
	accumulator  int64         // accumulator for timer.key
	currentIdx   int           // index of the current buckets
	started      int32
	mtx          sync.Mutex
}

type bucket struct {
	timers map[int64]*WheelTimer
}

// NewTimeWheelMgr
func NewTimeWheelMgr(interval time.Duration, bucketNum int) *TimeWheelMgr {
	if bucketNum < defaultBucketNum {
		bucketNum = defaultBucketNum
	}

	if interval < defaultInterval {
		interval = defaultInterval
	}

	mgr := &TimeWheelMgr{
		interval:     interval,
		bucketNum:    bucketNum,
		buckets:      make([]*bucket, bucketNum),
		timer2bucket: map[int64]int{},
		accumulator:  0,
	}

	for i := 0; i < bucketNum; i++ {
		mgr.buckets[i] = &bucket{timers: map[int64]*WheelTimer{}}
	}

	return mgr
}

// addTimer add timer to manager
func (mgr *TimeWheelMgr) addTimer(t *WheelTimer) {
	if t.key == 0 {
		key := atomic.AddInt64(&mgr.accumulator, 1)
		t.key = key
		t.mgr = mgr
	}

	delay := t.rt.expire.Sub(time.Now())
	if delay < mgr.interval {
		// execute directly if the expiration time is less than the manager.interval
		t.do()
	} else {
		mgr.mtx.Lock()
		defer mgr.mtx.Unlock()

		circle := int(int64(delay) / int64(mgr.interval) / int64(mgr.bucketNum))
		idx := (mgr.currentIdx + int(int64(delay)/int64(mgr.interval))) % mgr.bucketNum
		t.circle = circle

		mgr.timer2bucket[t.key] = idx
		mgr.buckets[idx].timers[t.key] = t
	}
}

// removeTimer remove timer from manager
func (mgr *TimeWheelMgr) removeTimer(t *WheelTimer) {
	mgr.mtx.Lock()
	defer mgr.mtx.Unlock()
	if idx, ok := mgr.timer2bucket[t.key]; ok {
		delete(mgr.timer2bucket, t.key)
		b := mgr.buckets[idx]
		delete(b.timers, t.key)
	}
}

func (mgr *TimeWheelMgr) run() {
	ticker := time.NewTicker(mgr.interval)
	for {
		<-ticker.C

		mgr.mtx.Lock()
		b := mgr.buckets[mgr.currentIdx]
		mgr.currentIdx = (mgr.currentIdx + 1) % mgr.bucketNum
		mgr.mtx.Unlock()
		for _, t := range b.timers {
			if t.circle > 0 {
				t.circle--
			} else {
				mgr.removeTimer(t)
				t.do()
			}
		}
	}
}

// registerTimer add timer to manager and run the manager first time
func (mgr *TimeWheelMgr) registerTimer(d time.Duration, repeated bool, f func()) *WheelTimer {
	if atomic.CompareAndSwapInt32(&mgr.started, 0, 1) {
		go mgr.run()
	}
	t := newWheelTimer(d, repeated, f)
	mgr.addTimer(t)
	return t
}

// OnceTimer waits for the duration to elapse and then calls f
// in its own goroutine. It returns a Timer. It's done once
func (mgr *TimeWheelMgr) OnceTimer(d time.Duration, f func()) Timer {
	return mgr.registerTimer(d, false, f)
}

// RepeatTimer waits for the duration to elapse and then calls f
// in its own goroutine. It returns a Timer. It can be used to
// cancel the call using its Stop method.
func (mgr *TimeWheelMgr) RepeatTimer(d time.Duration, f func()) Timer {
	return mgr.registerTimer(d, true, f)
}

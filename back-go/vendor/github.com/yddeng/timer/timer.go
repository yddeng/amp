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
	"fmt"
	"log"
	"runtime"
	"time"
)

type Timer interface {
	// Reset changes the timer to expire after duration d.
	// It returns true if the timer had been active or had executed,
	// false if the timer been stopped.
	Reset(duration time.Duration) bool

	// Stop prevents the Timer from firing.
	// It returns true if the call stops the timer, false if the timer has already
	// been stopped.
	Stop() bool
}

type runtimeTimer struct {
	expire   time.Time     // expire time
	fn       func()        // a function will be executed when expired
	repeated bool          // repeat task
	period   time.Duration // period of timer, it active if repeated is true
	stopped  int32         // timer is stopped
}

func goFunc(fn func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				buf := make([]byte, 2048)
				l := runtime.Stack(buf, false)
				log.Printf(fmt.Sprintf("%v: %s", r, buf[:l]))
			}
		}()
		fn()
	}()
}

func sendSignal(ch chan struct{}) {
	select {
	case ch <- struct{}{}:
	default:
	}
}

// when is a helper function for setting the 'expire' field of a runtimeTimer.
// It returns what the time will be, in time.Time, Duration d in the future.
func when(d time.Duration) time.Time {
	if d <= 0 {
		return time.Now()
	}
	return time.Now().Add(d)
}

// TimerMgr interface
type TimerMgr interface {
	// OnceTimer waits for the duration to elapse and then calls f
	// in its own goroutine. It returns a Timer. It's done once
	OnceTimer(d time.Duration, f func()) Timer

	// RepeatTimer waits for the duration to elapse and then calls f
	// in its own goroutine. It returns a Timer. It can be used to
	// cancel the call using its Stop method.
	RepeatTimer(d time.Duration, f func()) Timer
}

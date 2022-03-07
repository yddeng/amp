# timer
timer(heap_timer, timing_wheel) in Golang

```
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
```

```
type TimerMgr interface {
	// OnceTimer waits for the duration to elapse and then calls f
	// in its own goroutine. It returns a Timer. It's done once
	OnceTimer(d time.Duration, f func()) Timer

	// RepeatTimer waits for the duration to elapse and then calls f
	// in its own goroutine. It returns a Timer. It can be used to
	// cancel the call using its Stop method.
	RepeatTimer(d time.Duration, f func()) Timer
}
```

## heap timer

精度较高。

## timing wheel

低精度 timer, 最低精度为毫秒。

如果时间轮精度为10ms， 那么他的误差在 （0，10）ms之间。如果一个任务延迟 500ms，那它的执行时间在490～500ms之间。
出错的概率均等的情况下，那么这个出错可能会延迟或提前最小刻度的一半，在这里就是10ms/2=5ms.

故，时间轮的`interval` 在总延迟时间上，应该不足以影响延迟执行函数处理的事务。

时间轮的一个周期为 `interval`*`bucketNum`
```
NewTimeWheelMgr(interval time.Duration, bucketNum int)
```

## example

```
func TestNewHeapTimerMgr(t *testing.T) {
	mgr := NewHeapTimerMgr()

	fmt.Println("new---", time.Now().String())
	timer1 := mgr.OnceTimer(time.Second, func() {
		fmt.Println("once1", time.Now().String())
	})

	timer2 := mgr.RepeatTimer(time.Second*2, func() {
		fmt.Println("repeat1", time.Now().String())
	})

	// 立即执行
	mgr.OnceTimer(0, func() {
		fmt.Println("once3", time.Now().String())
	})
	
	//mgr.RepeatTimer(0, func() {
	//	fmt.Println("repeat4", time.Now().String())
	//})

	go func() {
		time.Sleep(time.Second * 5)
		fmt.Println("reset---", time.Now().String())
		fmt.Println(timer1.Reset(time.Second * 3))
		fmt.Println(timer2.Reset(time.Second))
	}()

	go func() {
		time.Sleep(time.Second * 10)
		fmt.Println("stop---", time.Now().String())
		timer1.Stop()
		timer2.Stop()
		fmt.Println(timer1.Reset(time.Second))
	}()

	time.Sleep(time.Second * 20)
}
```


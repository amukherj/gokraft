package once

import (
	"sync/atomic"
	"time"
)

// Only one of a set of competing goroutines complete
// a task. The others wait for it to complete.
type UntilSuccess struct {
	occupied int32
	done     int32
	sleepFor time.Duration
}

func NewUntilSuccess(sleepFor time.Duration) *UntilSuccess {
	return &UntilSuccess{
		occupied: 0,
		done:     0,
		sleepFor: sleepFor,
	}
}

func (u *UntilSuccess) Do(task func() error) {
	for atomic.LoadInt32(&u.done) == 0 {
		if atomic.CompareAndSwapInt32(&u.occupied, 0, 1) {
			err := task()
			if err == nil {
				atomic.StoreInt32(&u.done, 1)
				continue
			}
			atomic.StoreInt32(&u.occupied, 0)
		}
		time.Sleep(u.sleepFor)
	}
}

// Run at most one instance of a task at any point in time.
// A new goroutine is initiated only if there is none active.
type UniqueTask struct {
	occupied int32
	active   int32
	sleepFor time.Duration
}

func NewUniqueTask(sleepFor time.Duration) *UniqueTask {
	return &UniqueTask{
		occupied: 0,
		active:   0,
		sleepFor: sleepFor,
	}
}

// Callback called by the task on completion.
func (ua *UniqueTask) Done() {
	atomic.StoreInt32(&ua.active, 0)
}

func (ua *UniqueTask) DoAsync(task func(*UniqueTask)) bool {
	for atomic.LoadInt32(&ua.active) == 0 {
		if atomic.CompareAndSwapInt32(&ua.occupied, 0, 1) {
			go task(ua)
			atomic.StoreInt32(&ua.active, 1)
			return true
		}
		time.Sleep(ua.sleepFor)
	}
	return false
}

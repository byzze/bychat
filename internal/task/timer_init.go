package task

import "time"

// TimerFunc 定时器执行函数
type TimerFunc func(interface{}) bool

// Timer 定时器
func Timer(delay, tick time.Duration, fun TimerFunc, param interface{}, funcDefer TimerFunc, paramDefer interface{}) {
	go func() {
		defer func() {
			if funcDefer != nil {
				funcDefer(paramDefer)
			}
		}()

		if fun == nil {
			return
		}
		// 初始delay开始执行
		t := time.NewTimer(delay)
		defer t.Stop()

		for {
			select {
			case <-t.C:
				// 退出定时器
				if fun(param) == false {
					return
				}
				// 后续多少time开始执行
				t.Reset(tick)
			}
		}
	}()
}

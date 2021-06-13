package main

import (
	"context"
	"fmt"
	"strconv"
	"time"
)

func main() {

	t := tracker{ch: make(chan string, 10), stop: make(chan struct{})}
	// tracker 启动 goroutine 监控任务
	go t.run()

	// 异步生成20个任务
	for i := 0; i < 10; i++ {
		go t.event(context.Background(), "event"+strconv.Itoa(i))
	}
	time.Sleep(time.Duration(5) * time.Second)

	// 期望程序在 3s 内 shutdown
	context, cancelFunc := context.WithTimeout(context.Background(), time.Duration(3)*time.Second)
	defer func() {
		cancelFunc()
		fmt.Println("-- end to shutdown--")
	}()
	fmt.Println("-- start to shutdown--")
	t.shutdown(context)
}

type tracker struct {
	ch   chan string
	stop chan struct{}
}

/**
event: event1
event: event6
event: event7
event: event5
event: event3
event: event2
event: event8
event: event4
event: event0
event: event9
run: event6
run: event3
run: event4
run: event9
-- start to shutdown--
run: event7
run: event8
run: event1
shutdown timeout
-- end to shutdown--
*/
func (t *tracker) event(context context.Context, data string) error {
	select {
	case t.ch <- data:
		fmt.Println("event:", data)
		return nil
	case <-context.Done():
		return context.Err()
	}
}

func (t *tracker) run() {
	for c := range t.ch { // 2.1. 结束 range
		time.Sleep(time.Duration(1) * time.Second)
		//time.Sleep(time.Duration(5) * time.Second) // 2.2. 结束 range
		fmt.Println("run:", c)
	}
	t.stop <- struct{}{} // 3. stop 接收数据
}

func (t *tracker) shutdown(context context.Context) error {
	close(t.ch) // 1. close ch
	select {
	case <-t.stop: // 4.1. 在期望时间内 shutdown
		fmt.Println("shutdown success")
		return nil
	case <-context.Done(): // 4.2. 超时
		fmt.Println("shutdown timeout")
		return context.Err()
	}
}

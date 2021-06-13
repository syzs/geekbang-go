package main

import (
	"fmt"
	"sync"
)

/**
➜  go build -race race.go
➜ ./race
==================
WARNING: DATA RACE
Write at 0x000001174b48 by goroutine 7:
  main.work()
      /Users/shaoyongzhen/code/goup/geekbang/geekbang-go/src/concurrency/basic/race.go:23 +0x74

Previous read at 0x000001174b48 by goroutine 6:
  main.work()
      /Users/shaoyongzhen/code/goup/geekbang/geekbang-go/src/concurrency/basic/race.go:20 +0x47

Goroutine 7 (running) created at:
  main.main()
      /Users/shaoyongzhen/code/goup/geekbang/geekbang-go/src/concurrency/basic/race.go:14 +0x75

Goroutine 6 (finished) created at:
  main.main()
      /Users/shaoyongzhen/code/goup/geekbang/geekbang-go/src/concurrency/basic/race.go:14 +0x75
==================
Found 1 data race(s)

*/
var wg sync.WaitGroup
var count int

func main() {
	race()
}
func race() {
	for i := 1; i <= 2; i++ {
		wg.Add(1)
		go work(i)
	}
	wg.Wait()
	fmt.Println("final count:", count)
}
func work(id int) {
	for i := 0; i < 2; i++ {
		value := count
		//time.Sleep(time.Nanosecond)
		value++
		count = value
	}
	wg.Done()
}

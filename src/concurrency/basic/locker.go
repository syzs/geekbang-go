package main

import (
	"fmt"
	"sync"
)

/**
// A Locker represents an object that can be locked and unlocked.
type Locker interface {
   Lock()
   Unlock()
}
 */

/**
// Lock locks m.
// If the lock is already in use, the calling goroutine
// blocks until the mutex is available.
func (m *Mutex) Lock() {}

// Unlock unlocks m.
// It is a run-time error if m is not locked on entry to Unlock.
//
// A locked Mutex is not associated with a particular goroutine.
// It is allowed for one goroutine to lock a Mutex and then
// arrange for another goroutine to unlock it.
func (m *Mutex) Unlock() {}
 */
var locker sync.Mutex

func main() {
	go rrecover()
	//go rrecover()
	for {
		continue
	}
}

/**
fatal error: all goroutines are asleep - deadlock!
注意，这种由 Go 语言运行时系统自行抛出的 panic 都属于致命错误，都是无法被恢复的，
调用recover函数对它们起不到任何作用。也就是说，一旦产生死锁，程序必然崩溃。
 */
func rrecover() {
	defer func() {
		if p := recover(); p != nil {
			fmt.Println(p)
		}
	}()
	locker.Lock()
	fmt.Println("lock success")
	locker.Unlock()
	locker.Unlock()
}

package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"sync/atomic"
	"time"
)

/**
* 互斥锁的使用，可以保证临界区的代码同时只被一个goroutine执行，但不能保证执行的过程中不被中断，即CPU被其他goroutine获取 == 互斥锁虽然可以保证临界区中代码的串行执行，但却不能保证这些代码执行的原子性（atomicity）。
* 优点：原子操作在执行的过程中是不允许中断的。在底层，由CPU提供芯片级别的支持，所有绝对有效。这使得原子操作可以完全消除竞态条件，并能够保证并发的安全性。并且执行速度通常会比其他同步工具高出好几个数量级。
* 缺点：原子操作不能被中断，所以需要足够简单，足够快。一个需要长时间处理的不可中断的操作，会对计算机的执行指令的效率产生莫大的影响。因为，操作系统层面只对二进制位或整数的原子操作提供了支持。
* Go的原子操作是基于CPU和操作系统的，只对少量数据类型提供了原子操作。
* sync/atomic包中的函数可以做的原子操作有：加法（Add）、比较并交换（CompareAndSwap）、加载（Load）、存储（Store）和交换（Swap）。
* sync/atomic包中的函数支持的数据类型有：int32、int64、uint32、uint64、uintptr，以及unsafe包中的Pointer。不过，针对unsafe.Pointer类型，该包并未提供进行原子加法操作的函数。
* 传入参数都是指针类型：因为原子操作函数需要的是被操作值的指针，而不是这个值本身；被传入函数的参数值都会被复制，像这种基本类型的值一旦被传入函数，就已经与函数外的那个值毫无关系了。
 */
func main() {

	//spin()

	//forAndCAS2()

	value()
}

/**
The number: 0
The number: 2
The number: 4
The number: 6
The number: 8
The number: 10
The number has gone to zero.
*/
func spin() {
	sign := make(chan struct{}, 2)
	num := int32(0)
	fmt.Printf("The number: %d\n", num)
	go func() { // 定时增加num的值。
		defer func() {
			sign <- struct{}{}
		}()
		for {
			time.Sleep(time.Millisecond * 500)
			newNum := atomic.AddInt32(&num, 2)
			fmt.Printf("The number: %d\n", newNum)
			if newNum == 10 {
				break
			}
		}
	}()
	go func() { // 定时检查num的值，如果等于10就将其归零。
		defer func() {
			sign <- struct{}{}
		}()
		for {
			if atomic.CompareAndSwapInt32(&num, 10, 0) {
				fmt.Println("The number has gone to zero.")
				break
			}
			time.Sleep(time.Millisecond * 500)
		}
	}()
	<-sign
	<-sign
}

/**
[1-0] The number: 2
[2-0] The CAS number: 2 failed.
[1-1] The number: 4
[2-1] The CAS number: 4 failed.
[2-2] The number: 6
[1-2] The CAS number: 6 failed.
[2-3] The number: 8
[1-3] The CAS number: 8 failed.
[1-4] The number: 10
[2-4] The CAS number: 10 failed.
[2-5] The number: 12
[1-5] The CAS number: 12 failed.
[1-6] The number: 14
[2-6] The CAS number: 14 failed.
[2-7] The number: 16
[1-7] The CAS number: 16 failed.
[1-8] The number: 18
[2-8] The CAS number: 18 failed.
[2-9] The number: 20
[1-9] The CAS number: 20 failed.
*/
func forAndCAS2() {
	sign := make(chan struct{}, 2)
	num := int32(0)
	max := int32(20)
	for _, i := range []int32{1, 2} {
		go func(id, max int32) {
			defer func() {
				sign <- struct{}{}
			}()
			for i := 0; ; i++ {
				currNum := atomic.LoadInt32(&num)
				if currNum >= max {
					break
				}
				time.Sleep(time.Millisecond * 200)
				newNum := currNum + 2
				if atomic.CompareAndSwapInt32(&num, currNum, newNum) {
					fmt.Printf("[%d-%d] The number: %d \n", id, i, newNum)
				} else {
					fmt.Printf("[%d-%d] The CAS number: %d failed. \n", id, i, newNum)
				}
			}
		}(i, max)
	}
	for i := 0; i < cap(sign); i++ {
		<-sign
	}
}

/**
// A Value provides an atomic load and store of a consistently typed value.
// Values can be created as part of other data structures.
// The zero value for a Value returns nil from Load.
// Once Store has been called, a Value must not be copied.
//
// A Value must not be copied after first use.
type Value struct {
   noCopy noCopy
   v interface{}
}
 */

func value() {
	// 示例1。
	var box atomic.Value
	fmt.Println("Copy box to box2.")
	box2 := box // 原子值在真正使用前可以被复制。
	v1 := [...]int{1, 2, 3}
	fmt.Printf("Store %v to box.\n", v1)
	box.Store(v1)
	fmt.Printf("The value load from box is %v.\n", box.Load())
	fmt.Printf("The value load from box2 is %v.\n", box2.Load())
	fmt.Println("----------------------------")

	// 示例2。
	v2 := "123"
	fmt.Printf("Store %s to box2.\n", v2)
	box2.Store(v2) // 这里并不会引发panic。
	fmt.Printf("The value load from box is %v.\n", box.Load())
	fmt.Printf("The value load from box2 is %q.\n", box2.Load())
	fmt.Println("----------------------------")

	// 示例3。
	fmt.Println("Copy box to box3.")
	box3 := box // 原子值在真正使用后不应该被复制！
	fmt.Printf("The value load from box3 is %v.\n", box3.Load())
	v3 := 123
	fmt.Printf("Store %d to box3.\n", v3)
	//box3.Store(v3) // 这里会引发一个panic，报告存储值的类型不一致。panic: sync/atomic: store of inconsistently typed value into Value
	_ = box3
	fmt.Println("----------------------------")

	// 示例4。
	var box4 atomic.Value
	v4 := errors.New("something wrong")
	fmt.Printf("Store an error with message %q to box4.\n", v4)
	box4.Store(v4)
	v41 := io.EOF
	fmt.Println("Store a value of the same type to box4.")
	box4.Store(v41)
	v42, ok := interface{}(&os.PathError{}).(error)
	if ok {
		fmt.Printf("Store a value of type %T that implements error interface to box4.\n", v42)
		//box4.Store(v42) // 这里会引发一个panic，报告存储值的类型不一致。panic: sync/atomic: store of inconsistently typed value into Value
	}
	fmt.Println("----------------------------")
}
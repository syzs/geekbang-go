package main

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"math/rand"
	"time"
)

/**
1. 每一个case表达式，必须包含一个代表发送或接受的表达式。当包含多个表达式时，从左到右被顺序求值；
2. 多个case表达式，判断是否求值成功的顺序是从上到下；
3. 当一个case表达式在被求值时，相应的操作处于阻塞状态，即求值不成功，认为当前case不满足条件；
4. 当所有case中的表达式都计算完成后，才会开始选择分支。如果没有命中的分支，select会阻塞，一旦有分支可以被命中，select所在的goroutine会被唤醒，选中的分支就会开始执行；
5. 当有多个分支满足条件，select会随机选择一条进行执行；
6. select只能有一条默认分支，只有在其他case都没有命中的情况下才会执行，与分支所在的位置无关；
7. select语句的执行及case表达式的计算，都是独立的，是并发在计算的。
*/

func main() {
	//example2()

	//example301()

	//example302()

	//example4()

	example5()

	//blockSelect()
}

// 空的 select 会永远阻塞
func blockSelect() {
	go func() {
		fmt.Println("block select")
	}()
	// fatal error: all goroutines are asleep - deadlock!
	select {}
}

// 示例1。
func example1() {
	// 准备好几个通道。
	intChannels := [3]chan int{
		make(chan int, 1),
		make(chan int, 1),
		make(chan int, 1),
	}
	// 所有分支都会阻塞
	select {
	case elem := <-intChannels[0]:
		fmt.Printf("The first candidate case is selected, the element is %d.\n", elem)
	case elem := <-intChannels[1]:
		fmt.Printf("The second candidate case is selected, the element is %d.\n", elem)
	case elem := <-intChannels[2]:
		fmt.Printf("The third candidate case is selected, the element is %d.\n", elem)
	}
}

// 示例2。
func example2() {
	// 准备好几个通道。
	intChannels := [3]chan int{
		make(chan int, 1),
		make(chan int, 1),
		make(chan int, 1),
	}

	// 随机选择一个通道，并向它发送元素值。
	index := rand.Intn(3)
	fmt.Printf("The index: %d\n", index)
	intChannels[index] <- index // 向其中一个通道中发送了数据

	// 哪一个通道中有可取的元素值，哪个对应的分支就会被执行。
	// 只要有一个分支命中，则不会阻塞整个线程
	select {
	case elem := <-intChannels[0]:
		fmt.Printf("The first candidate case is selected, the element is %d.\n", elem)
	case elem := <-intChannels[1]:
		fmt.Printf("The second candidate case is selected, the element is %d.\n", elem)
	case elem := <-intChannels[2]:
		fmt.Printf("The third candidate case is selected, the element is %d.\n", elem)
	}
}

// 示例3。
func example301() {
	// 准备好几个通道。
	intChannels := [3]chan int{
		make(chan int, 1),
		make(chan int, 1),
		make(chan int, 1),
	}

	// 哪一个通道中有可取的元素值，哪个对应的分支就会被执行。
	// 只要有一个分支命中，则不会阻塞整个线程
	select {
	case elem := <-intChannels[0]:
		fmt.Printf("The first candidate case is selected, the element is %d.\n", elem)
	case elem := <-intChannels[1]:
		fmt.Printf("The second candidate case is selected, the element is %d.\n", elem)
	case elem := <-intChannels[2]:
		fmt.Printf("The third candidate case is selected, the element is %d.\n", elem)
	default: // case的分支中没有命中的，只要有default分支就会命中default分支
		fmt.Println("No candidate case is selected!")
	}
}

var channels = [3]chan int{
	nil,
	make(chan int),
	nil,
}

var numbers = []int{1, 2, 3}

func example302() {
	select {
	case getChan(0) <- getNumber(0): // nil通道 阻塞
		fmt.Println("The first candidate case is selected.")
	case getChan(1) <- getNumber(1): // 非缓冲通道,发送数据后,会阻塞
		fmt.Println("The second candidate case is selected.")
	case getChan(2) <- getNumber(2): // nil通道 阻塞
		fmt.Println("The third candidate case is selected")
	default: // case 分支全部阻塞，必然执行到default分支
		fmt.Println("No candidate case is selected!")
	}
}

func getNumber(i int) int {
	fmt.Printf("numbers[%d]\n", i)
	return numbers[i]
}

func getChan(i int) chan int {
	fmt.Printf("channels[%d]\n", i)
	return channels[i]
}

// 如果在select语句中发现某个通道已关闭，那么应该怎样屏蔽掉它所在的分支？
func example4() {
	ch1 := make(chan int, 1)
	ch2 := make(chan int, 1)

	ch1 <- 1
	ch2 <- 2
	close(ch1)

	for i := 0; i < 5; i++ {
		select {
		case _, ok := <-ch1: // 1.接收数据，ch1 // 3.关闭通道 close // 4.通道变成nil或无缓冲的通道，分支阻塞
			if ok {
				fmt.Println("ch1")
			} else {
				fmt.Println("ch1 close")
				ch1 = nil // or ch1 = make(chan int)
				break
			}
		case <-ch2: // 2.接收数据，ch2
			fmt.Println("ch2")
		default: // 4、5.case分支全部阻塞，执行默认分支
			fmt.Println("no case selected")
		}
	}
}

func example5() (err error) {

	defer func() {
		if err != nil {
			fmt.Printf("%+v", err)
		}
	}()

	cxt, _ := context.WithTimeout(context.Background(), time.Duration(20)*time.Nanosecond)

	type result struct {
		record string
		err    error
	}

	//ch := make(chan result) // 无缓冲的channel，写入之后没有读出，线程阻塞，会出现泄漏的情况
	ch := make(chan result, 1) // 设置缓冲容量为1的channel,及时操作超时没有读操作，操作线程不会阻塞
	go func() {
		record, err := search("term")
		ch <- result{record, err}
	}()

	select {
	case <-cxt.Done(): // 超时
		return errors.New("timeout")
	case result := <-ch:
		if result.err != nil {
			return errors.Wrap(result.err, "")
		}
		fmt.Println("result:", result.record)
		return nil
	}
}

// 模拟一个耗时操作
func search(term string) (string, error) {
	time.Sleep(100 * time.Nanosecond)
	return term, nil
}

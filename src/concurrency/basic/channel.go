package main

/**
可以利用channel在多个goroutine之间通信，所以channel是线程安全的。Go自带的、唯一的线程安全的数据类型。
一个channel相当于一个FIFO的队列，底层数据是环形链表
1. 对于同一个通道，发送操作之间只互斥的，接收操作也是互斥的；
    1. Go语言的运行时系统只会执行同一通道的任意个发送操作中的某一个，直到元素完全被复制进通道中，才会执行其他的发送操作；接收操作同上，直到元素值完全被移除通道；
    2. 对于同一个元素，发送操作和接收操作也是互斥的，未完全复制进通道的值不会被接收；
    3. 元素进入通道是被复制（浅层复制）的，即进入通道的并不是元素值本身，而是其副本。
    4. 元素值在被接收时，首先生成元素值的副本并准备给接收方，然后从通道中删除这个元素。
2. 发送操作和接收操作对元素值的操作都是不可分割的；
    1. 发送操作要么复制进了元素进通道，要么没复制元素进通道，不会出现复制部分的情况；
    2. 接收操作在准备好元素值的副本之后，一定会删除掉通道中的原值，绝不会出现通道中仍有残留的情况；
    3. 保证元素值的完整性，也保证操作的完整性。
3. 发送操作在完全完成之前会被阻塞；接收操作也是如此。
    1. 发送操作包括复制元素副本及将副本放入通道中，在完成发送操作前，发起发送操作的代码语句会阻塞，直到发送操作完成，运行时系统会通知发起请求的goroutine去争取系统资源，继续往下执行；
    2. 接收操作包括复制通道内元素值、放置副本到接收方、删除通道内元素，在完成接收操作前，同上；
    3. 实现操作的互斥和元素的完整性。
*/
import "fmt"

func main() {

	//channelCap()

	//channelBlock()

	//copy()

	senderReceiver()

	//rangeChannel()
}
func rangeChannel() {
	ch := make(chan int, 10)
	go func() {
		for i := 1; i <= 10; i++ {
			ch <- i
		}
		close(ch)
	}()

	for i := 0; i < cap(ch); i++ {
		fmt.Println(<-ch)
	}

	// clone ch，才能正常退出循环，否则 fatal error: all goroutines are asleep - deadlock!
	for c := range ch { //
		fmt.Println(c)
	}
}
func channelCap() {
	ch1 := make(chan int, 3) // 初始化一个缓冲容量为3，数据类型为int的通道
	ch1 <- 2                 // 发送数据
	ch1 <- 1
	ch1 <- 3
	elem1 := <-ch1 // 接收数据
	fmt.Printf("The first element received from channel ch1: %v\n", elem1)
	elem1 = <-ch1 // 接收数据
	fmt.Printf("The second element received from channel ch1: %v\n", elem1)
	elem1 = <-ch1 // 接收数据
	fmt.Printf("The third element received from channel ch1: %v\n", elem1)

	close(ch1)
	// fatal error: all goroutines are asleep - deadlock!
	//elem1 = <-ch1 // 接收数据
}

func channelBlock() {

	// 未初始化的 channel，写入数据操作会死锁
	//var ch  chan struct{}
	//ch <- struct {}{}

	//ch0 := make(chan int) // 无缓冲的通道
	////ch0 <- 1 // 无缓冲容量，数据没有被接收前阻塞当前线程，造成永久阻塞
	//
	//go func() {
	//	fmt.Println(<-ch0)
	//}()
	//ch0 <- 1
	//
	//// 示例1。
	//ch1 := make(chan int, 1)
	//ch1 <- 1
	////ch1 <- 2 // 通道已满，因此这里会造成阻塞。
	//
	//// 示例2。
	//ch2 := make(chan int, 1)
	////elem, ok := <-ch2 // 通道已空，因此这里会造成阻塞。
	////_, _ = elem, ok
	//ch2 <- 1
	//
	//// 示例3。
	//var ch3 chan int
	////ch3 <- 1 // 通道的值为nil，因此这里会造成永久的阻塞！
	////<-ch3 // 通道的值为nil，因此这里会造成永久的阻塞！
	//_ = ch3
	//
	//close(ch0)
	//close(ch1)
	//close(ch2)
	////close(ch3) // panic: close of nil channel

}

func copy() {
	m := map[int]int{1: 1}
	ch1 := make(chan map[int]int, 1)
	ch1 <- m
	m[2] = 2
	elem1 := <-ch1
	fmt.Printf("ch1: %v\n", elem1) // ch1: map[1:1 2:2] 引用类型

	a := []int{1, 2}
	ch2 := make(chan []int, 1)
	ch2 <- a
	a[1] = 100
	elem2 := <-ch2
	fmt.Printf("ch2: %v, &ch2:%p a:%v, &a:%p \n", elem2, elem2, a, a) // ch2: [1 100], &ch2:0xc00012e040 a:[1 100], &a:0xc00012e040

	ch2 <- a
	a = append(a, 3) // 扩容后，指向了的新的底层数组
	elem2 = <-ch2
	fmt.Printf("ch2: %v, &ch2:%p a:%v, &a:%p \n", elem2, elem2, a, a) // ch2: [1 100], &ch2:0xc0000a6040 a:[1 100 3], &a:0xc0000ac020

	close(ch1)
	close(ch2)
}

func senderReceiver() {
	ch := make(chan int)

	go func() {
		for i := 1; i <= 10; i++ {
			fmt.Printf("--Sender send: %d --\n", i)
			ch <- i
		}
		fmt.Println("Sender close channel")
		close(ch) // 关闭通道
	}()

	for {
		v, ok := <-ch
		if ok {
			fmt.Println("**Receiver receive:", v)
		} else {
			fmt.Println("**Receiver close channel")
			break
		}
	}

}

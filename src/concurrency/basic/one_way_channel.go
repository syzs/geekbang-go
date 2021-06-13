package main

import (
	"fmt"
	"math/rand"
)

func main() {
	// 示例1。
	// 只能发不能收的通道。
	var uselessChan = make(chan<- int, 1)
	// 只能收不能发的通道。
	var anotherUselessChan = make(<-chan int, 1)
	// 这里打印的是可以分别代表两个通道的指针的16进制表示。
	fmt.Printf("The useless channels: %v, %v\n", uselessChan, anotherUselessChan)

	// 示例2。
	intChan1 := make(chan int, 3)
	SendInt(intChan1) // Go自动把双向通道转成了函数所需的单向通道
	close(intChan1)

	// 示例4。
	intChan2 := getIntChan()
	for elem := range intChan2 {
		fmt.Printf("The element in intChan2: %v\n", elem)
	}

	// 示例5。
	_ = GetIntChan(getIntChan)
}

// 示例2。
func SendInt(ch chan<- int) { // 在这个函数中的代码只能向参数ch发送元素值，而不能从它那里接收元素值。这就起到了约束函数行为的作用。
	ch <- rand.Intn(1000)
}

// 示例3。
type Notifier interface {
	SendInt(ch chan<- int) // 在接口中对实现做出了约束
}

// 示例4。
func getIntChan() <-chan int { // 函数getIntChan会返回一个<-chan int类型的通道，这就意味着得到该通道的程序，只能从通道中接收元素值。这实际上就是对函数调用方的一种约束了。
	num := 5
	ch := make(chan int, num)
	for i := 0; i < num; i++ {
		ch <- i
	}
	close(ch)
	return ch
}

// 示例5。
type GetIntChan func() <-chan int

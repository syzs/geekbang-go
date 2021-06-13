package main

/**
* Go 1.7后支持
* 一种非常通用的同步工具，可以提供一类代表上下文的值，是并发安全的
* 它的值不但可以被任意的扩散，还可以被用来传递额外的信息和信号
* 是可以繁衍的，可以通过一个Context值产生出任意个子值，这些子值可以携带起父值的属性和数据，也可以响应通过其父值传达的信号
* 所有的Context值构成代表上下文全貌的树形结构，树的作用是全局的。
* 树根（上下文根节点）是一个已经在contxt包中预定义好的Context值，是全局唯一的。可以通过 context.Background获取。跟节点不提供任何功能。

type Context interface {
   Deadline() (deadline time.Time, ok bool)
   Done() <-chan struct{} // 让调用方感知撤销当前Context值的信号。一旦当前的Context值被撤销，接受通道即会被关闭。对一个未包含任何元素值的通道来说，关闭会使任何针对它的接受操作立即结束
   Err() error // 得知撤销的原因，其值只可能等于context.Canceled(手动撤销)变量的值，或者context.DeadlineExceeded(超过到期时间导致的撤销)变量的值。
   Value(key interface{}) interface{}
}

Context类型的子值
* WithCancel：用于触发撤销信号的函数。撤销：终止程序对某种请求（如http请求）的响应，或取消对某种指令（如sql指令）的处理
* WithDeadline：产生一个会定时撤销的函数。通过内部的计时器来实现，到期进行自动撤销并释放掉其内部的计时器
* WithTimeout：同上
* WithValue：携带额外数据的函数

func WithCancel(parent Context) (ctx Context, cancel CancelFunc) {}

func WithDeadline(parent Context, d time.Time) (Context, CancelFunc) {}

func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc) {}

func WithValue(parent Context, key, val interface{}) Context {}
 */
import (
	"context"
	"fmt"
	"sync/atomic"
	"time"
)

func main() {

	//withValue()

	coordinateWithContext()
}

/**
撤销信号如何在上下文书传播
1. 撤销函数被调用后，对用的Context值会先关闭内部的接受通道，即Done方法会返回的那个通道
2. 向所有的子值传递撤销信号。子值继续向其子值传播。最后，这个Context值会断开与父值的关联
3. WithValue得到的Context子值不可被撤销。撤销信号遇到它们会直接跨过，并试图将撤销信号传递给它们的子值。
4. “撤销”这个操作是Context值能够协调多个 goroutine 的关键所在。撤销信号总是会沿着上下文树叶子节点的方向传播开来。
 */
func coordinateWithContext() {
	total := 12
	var num int32
	fmt.Printf("The number: %d [with context.Context]\n", num)
	cxt, cancelFunc := context.WithCancel(context.Background()) // 1.产生一个可撤销的Context值及一个撤销函数
	for i := 1; i <= total; i++ {
		go addNum(&num, i, func() {
			if atomic.LoadInt32(&num) == int32(total) {
				cancelFunc() // 2.调用撤销函数，撤销信号立即被传达给Content值，并有Done方法的结果值接收
			}
		})
	}
	<-cxt.Done() // 3.结束通道的接收
	fmt.Println("End.")
}
// addNum 用于原子地增加一次numP所指的变量的值。
func addNum(numP *int32, id int, deferFunc func()) {
	defer func() {
		deferFunc()
	}()
	for i := 0; ; i++ {
		currNum := atomic.LoadInt32(numP)
		newNum := currNum + 1
		time.Sleep(time.Millisecond * 200)
		if atomic.CompareAndSwapInt32(numP, currNum, newNum) {
			fmt.Printf("The number: %d [%d-%d]\n", newNum, id, i)
			break
		} else {
			fmt.Printf("The CAS operation failed. [%d-%d]\n", id, i)
		}
	}
}

/**
Context携带数据
1. key必须是可判等的。从中获取数据时，根据给定的键来查找相应的值。Context值并不是用字典来存储键和值的，后两者只是被简单地存储在前者的相应的字段中。
2. Context类型的Value方法就是被用来获取数据的。调用含数据的Context值的Value方法时，它会先判断给定的键，是否与当前值中存储的键相等，如果相等就把该值中存储的值直接返回，
	否则就到其父值中继续查找。如果其父值中仍然未存储相等的键，那么该方法就会沿着上下文根节点的方向一路查找下去。
3. 除了含数据的Context值以外，其他几种Context值都是无法携带数据的。因此，Context值的Value方法在沿路查找的时候，会直接跨过那几种值。
 */
func withValue(){
	keys := []int{
		20,
		30,
		60,
		61,
	}
	values := []string{
		"value in node2",
		"value in node3",
		"value in node6",
		"value in node6Branch",
	}

	rootNode := context.Background()
	node1, cancelFunc1 := context.WithCancel(rootNode)
	defer cancelFunc1()


	// 示例1。
	node2 := context.WithValue(node1, keys[0], values[0])
	node3 := context.WithValue(node2, keys[1], values[1])
	fmt.Printf("The value of the key %v found in the node3: %v\n",
		keys[0], node3.Value(keys[0])) // 自身没有，向父值查找
	fmt.Printf("The value of the key %v found in the node3: %v\n",
		keys[1], node3.Value(keys[1])) // 自身有，不向父值查找
	fmt.Printf("The value of the key %v found in the node3: %v\n",
		keys[2], node3.Value(keys[2])) // 自身和父值都没有，返回nil
	fmt.Println()

	// 示例2。
	node5 := context.WithValue(node1, 5, 5)
	node6 := context.WithValue(node5, 6, 6)
	fmt.Printf("The value of the key %v found in the node5: %v\n",
		5, node5.Value(5)) // 自身没有，向父值查找
	fmt.Printf("The value of the key %v found in the node6: %v\n",
		5, node6.Value(5)) // 自身没有，向父值查找
	node6 = context.WithValue(node1, 6, 6)
	fmt.Printf("The value of the key %v found in the node6: %v\n",
		5, node6.Value(5)) // node5 不再是 node6 的父值

	// 示例2。
	node3 = context.WithValue(node2, keys[0], values[1])
	fmt.Printf("The value of the key %v found in the node2: %v\n",
		keys[0], node2.Value(keys[0])) // 子值的修改对父值没有影响
	fmt.Printf("The value of the key %v found in the node3: %v\n",
		keys[0], node3.Value(keys[0]))
	fmt.Println()

	// 示例3。
	node2 = context.WithValue(node1, keys[0], values[0])
	fmt.Printf("The value of the key %v found in the node2: %v\n",
		keys[0], node2.Value(keys[0])) // 父值的修改不会影响到子值
	fmt.Printf("The value of the key %v found in the node3: %v\n",
		keys[0], node3.Value(keys[0]))
}
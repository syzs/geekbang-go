package main

import "fmt"

type Tom struct {
	id   int
	name string
}

func (b *Tom) Hello() {
	fmt.Println("Tom says: my name is", b.name)
}

type Jerry struct {
	name string
}

func (j *Jerry) Hello() {
	fmt.Println("Jerry says: my name is", j.name)
}

/**
type interface struct{
    Type unitptr // points to the type of the interface implements
    Data unitptr // holds the data for the interface's receiver
}
*/
type Worker interface {
	Hello()
}

/**

会出现 Jerry says: my name Tom 或 Tom says: my name Jerry 的情况

如果两个类型的数据结构不同的话，内存布局不相同，会出现panic的情况
panic: runtime error: invalid memory address or nil pointer dereference
[signal SIGSEGV: segmentation violation code=0x1 addr=0x5 pc=0x1065253]


goroutine 1 [running]:
fmt.(*buffer).writeString(...)
        /usr/local/Cellar/go/1.16.3/libexec/src/fmt/print.go:82
fmt.(*fmt).padString(0xc00006b150, 0x5, 0x10e8718)
        /usr/local/Cellar/go/1.16.3/libexec/src/fmt/format.go:110 +0x8e
*/
func main() {
	Tom := &Tom{name: "Tom", id: 10}
	Jerry := &Jerry{"Jerry"}

	var worker Worker = Tom // 虽然Tom是个指针，但Tom赋值给worker并不是原子操作的，需要更新worker的2个字段
	var loop0, loop1 func()

	loop0 = func() {
		worker = Tom
		go loop1()
	}

	loop1 = func() {
		worker = Jerry
		go loop0()
	}

	go loop0()

	for {
		worker.Hello()
	}

}

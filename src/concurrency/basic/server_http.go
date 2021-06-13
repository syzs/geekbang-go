package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

/**
nerver start a goroutine without knowing when it will stop
when it will terminate
how to terninate
*/

// example01
func serverApp01() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintln(writer, "hello")
	})
	if err := http.ListenAndServe("0.0.0.0:8080", mux); err != nil {
		/**
		std.Output(2, fmt.Sprint(v...))
		os.Exit(1)
		*/
		// 整个进程直接退出，不会执行defer
		// 除 init 和 main 意外，不建议使用 Fatal
		log.Fatal(err)
	}
}

func serverDebug01() {
	if err := http.ListenAndServe("0.0.0.0:8001", http.DefaultServeMux); err != nil {
		log.Fatal(err)
	}
}

func main01() {
	go serverApp01()
	go serverDebug01()
	select {}
}

// example02
/**
server:app, err:listen tcp 0.0.0.0:6666: bind: address already in use
main error: listen tcp 0.0.0.0:6666: bind: address already in use
debug shutdown
app shutdown
server:debug, err:http: Server closed
main error: http: Server closed

debug 启动成功，并且有一个 阻塞在 <-stop 的 goroutine；app 启动因 debug 占用端口而启动失败，返回error，也有一个 阻塞在 <-stop 的 goroutine
done 接收到 serverApp02 返回的error，打印错误
!stopped == true，这是 stopped 为 true，并 close stop;
serverApp02 和 serverDebug02 的各自 阻塞在 <-stop 的 goroutine，因为 stop 被 close 掉了，而进入非阻塞状态
debug 启动成功，因 shutdonw 而返回了一个 Server close 的错误，done 接收到 serverDebug02 返回的error，打印错误
!stopped == false，退出 main
 */
func server(name, addr string, handler http.Handler, stop <-chan struct{}) error {
	server := http.Server{
		Addr:    addr,
		Handler: handler,
	}
	go func() {
		<-stop // wait for stop single
		fmt.Println(name, "shutdown")
		server.Shutdown(context.Background())
	}()
	err := server.ListenAndServe() // block
	fmt.Printf("server:%s, err:%v\n", name, err)
	return err
}

func serverApp02(stop <-chan struct{}) error {
	return server("app", "0.0.0.0:6666", http.DefaultServeMux, stop)
}

// 端口冲突，后获取到上下文的线程启动失败
func serverDebug02(stop <-chan struct{}) error {
	return server("debug", "0.0.0.0:6666", http.DefaultServeMux, stop)
}
func main() {
	done := make(chan error, 2)

	stop := make(chan struct{})

	go func() {
		done <- serverApp02(stop)
	}()
	go func() {
		done <- serverDebug02(stop)
	}()

	var stopped bool
	for i := 0; i < cap(done); i++ {
		if err := <-done; err != nil {
			fmt.Println("main error:",err)
		}
		if !stopped { // 有个终止标志，避免重复 close channel
			stopped = true
			// 必须是 close 才能同时让多个goroutine的 <- 停止阻塞；
			// 如果是 ch = nil 将永久阻塞；
			// 如果是 ch <- struct{}，将只有一个 goroutine 能接收到数据，进行 shutdonw
			close(stop)
		}
	}
}

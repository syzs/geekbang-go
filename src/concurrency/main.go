package main

import (
	"context"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"os"
	"os/signal"
)

func server(ctx context.Context, name, addr string, handler http.Handler) error {
	log.Println("start server:", name)
	server := http.Server{
		Addr:    addr,
		Handler: handler,
	}
	go func() {
		<-ctx.Done()
		server.Shutdown(ctx)
	}()
	err := server.ListenAndServe()
	return errors.Wrap(err, "server "+name)
}

func serverApp(ctx context.Context) error {
	return server(ctx, "app", "0.0.0.0:6666", http.DefaultServeMux)
}

func serverDebug(ctx context.Context) error {
	//return server(ctx, "debug", "0.0.0.0:7777", http.DefaultServeMux)
	return errors.New("somethind wrong")
}

type serverApplication func(ctx context.Context) error

func main() {

	serverList := []serverApplication{serverApp, serverDebug}

	parentCtx, cancelFunc := context.WithCancel(context.Background())
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, os.Kill)

	defer func() {
		if cancelFunc != nil {
			cancelFunc()
		}
		close(signalCh)
	}()
	// 利用 context 注销消息的传播路径，父 ctx 注销，会注销子 ctx
	go func() {
		select {
		case s, ok := <-signalCh:
			if ok {
				log.Printf("Got signal:%v\n", s)
				cancelFunc()
			}
		}
	}()

	eg, ctx := errgroup.WithContext(parentCtx)

	for i, _ := range serverList {
		server := serverList[i]
		eg.Go(func() (err error) {
			err = server(ctx)
			log.Println(err)
			return
		})
	}
	if err := eg.Wait(); err != nil {
		log.Println("server exits")
	}
}

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

func main() {

	eg, ctx := errgroup.WithContext(context.Background())

	mux := http.NewServeMux()

	server := &http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	httpStop := make(chan struct{}, 1)

	eg.Go(func() error {
		startHttpServer(mux, server, httpStop)
		return nil

	})

	eg.Go(func() error {

		select {
		case <-httpStop:
			stopHttpServer(server)
			return errors.New("shutdown from http")
		case <-ctx.Done():
			return nil
		}

	})

	eg.Go(func() error {

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGHUP)

		select {
		case s := <-quit:
			stopHttpServer(server)
			return errors.Errorf("shutdown from signal %s", s)
		case <-ctx.Done():
			return nil
		}

	})

	fmt.Printf("program exit: %s\n", eg.Wait())

}

func stopHttpServer(server *http.Server) {
	log.Printf("http server shuting down...")

	ctx, cancle := context.WithTimeout(context.Background(), 1*time.Second)

	defer cancle()
	err := server.Shutdown(ctx)

	if err != nil {
		log.Printf("http server shutdown error: %s", err)
	}
}


func startHttpServer(mux *http.ServeMux, server *http.Server, stop chan<- struct{}) error {
	mux.HandleFunc("/hello", SayHello)
	mux.HandleFunc("/shutdown", func(writer http.ResponseWriter, request *http.Request) {
		log.Printf("get shutdown signal from http")
		stop <- struct{}{}
	})
	return server.ListenAndServe()
}


func SayHello(rw http.ResponseWriter, req *http.Request) {
	rw.Write([]byte("Hello World!"))
}

package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jerbe/go-pool"
)

func main() {
	dispatcher := gowork.NewDispatcher(1000)

	go dispatcher.Run()

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGKILL)

	go func() {
		for {
			select {
			case <-time.Tick(time.Nanosecond):
				err := dispatcher.PushFunc(func() error {
					return nil
				})
				if err != nil {
					log.Println(err)
				}
			case <-sig:
				return
			}
		}
	}()

	log.Println("main 1")
	time.Sleep(time.Second)
	log.Println("main 2")
	dispatcher.Stop()
	log.Println("main 3")

	<-sig
}

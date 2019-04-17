package gowork

import (
	"fmt"
	"log"
	"time"
)

// Worker 工作者
type Worker interface {
	Start()
	Stop()
}

type workerChan chan jobChan

type worker struct {
	id string

	dispatcher Dispatcher

	jobChannel jobChan

	quit chan bool
}

func (w *worker) start() {
	defer func() {
		if err := recover(); err != nil {
			go w.start()
		}
	}()

	for {
		select {
		case j := <-w.jobChannel:
			if err := j.Execute(); err != nil {
				log.Println("worker.start() -> error", err)
			}
			w.register()

		case <-w.quit:
			fmt.Println(w.id, "worker quit")
			return
		}
	}
}

func (w *worker) register() {
	d := w.dispatcher.(*dispatcher)
	d.workPool <- w.jobChannel
}

// Start 工作者开始
func (w *worker) Start() {
	go func() {
		w.register()
		w.start()
	}()
}

// Stop 工作者结束
func (w *worker) Stop() {
	w.quit <- true
}

func newWorker(d Dispatcher) Worker {
	return &worker{
		id: fmt.Sprintf("%d", time.Now().UnixNano()),

		jobChannel: make(jobChan),

		dispatcher: Dispatcher(d),

		quit: make(chan bool),
	}
}

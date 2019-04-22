package gowork

import (
	"fmt"
	"sync"
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

	working bool

	workLocker sync.Mutex
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
				if ExecuteErrorHandle != nil {
					ExecuteErrorHandle(err)
				}
			}
			w.register()
		case <-w.quit:
			return
		}
	}
}

func (w *worker) register() {
	d := w.dispatcher.(*dispatcher)
	d.workPool <- w.jobChannel
}

// Run 工作者开始
func (w *worker) Start() {
	w.workLocker.Lock()
	w.working = true
	defer w.workLocker.Unlock()
	go func() {
		w.register()
		w.start()
	}()
}

// Stop 工作者结束
func (w *worker) Stop() {
	w.workLocker.Lock()
	if !w.working {
		return
	}
	defer w.workLocker.Unlock()
	w.quit <- true
}

func newWorker(d Dispatcher) Worker {
	return &worker{
		id: fmt.Sprintf("%d", time.Now().UnixNano()),

		jobChannel: make(jobChan),

		dispatcher: d,

		quit: make(chan bool),
	}
}

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
		w.register()
		select {
		case j := <-w.jobChannel:
			if err := j.Execute(); err != nil {
				if ExecuteErrorHandle != nil {
					ExecuteErrorHandle(err)
				}
			}

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
	defer w.workLocker.Unlock()
	w.working = true
	if w.working {
		return
	}
	go func() {
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

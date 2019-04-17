package gowork

import (
	"errors"
	"log"
)

// Dispatcher 分配器
type Dispatcher interface {
	Run() error
	Push(Job) error
	PushFunc(func() error) error
	Stop() error
}

// NewDispatcher 返回一个新的分配器
func NewDispatcher(cap int) Dispatcher {
	if cap <= 0 {
		panic(errors.New("cap can't not less than zero"))
	}

	workers := make([]Worker, cap, cap)

	workerPool := make(workerChan, cap)

	jobQueue := make(jobChan)

	quit := make(chan bool)

	d := &dispatcher{
		cap: cap,

		workers: workers,

		jobQueue: jobQueue,

		workPool: workerPool,

		quit: quit,
	}

	for i := 0; i < cap; i++ {
		w := newWorker(d)
		workers[i] = w
	}

	return d
}

type dispatcher struct {
	cap int

	workers []Worker

	jobQueue jobChan

	workPool workerChan

	quit chan bool

	running bool
}

func (d *dispatcher) run() {
	for {
		select {
		case job, ok := <-d.jobQueue:
			if ok {
				go func(j Job) {
					jc, ok := <-d.workPool
					if ok {
						jc <- j
					}
				}(job)
			}

		case <-d.quit:
			d.stop()
			return
		}
	}
}

// Run 开始并执行分发
func (d *dispatcher) Run() error {
	if d.running {
		return nil
	}

	d.running = true

	if d.workers == nil || len(d.workers) == 0 {
		return errors.New("worker is nil")
	}

	if d.jobQueue == nil {
		return errors.New("job queue was nil")
	}

	if d.workPool == nil {
		return errors.New("work pool was nil")
	}

	for i := 0; i < len(d.workers); i++ {
		d.workers[i].Start()
	}

	d.run()

	return nil
}

// Push 推入一条任务
func (d *dispatcher) Push(j Job) error {
	if !d.running {
		return nil
	}

	d.jobQueue <- j
	return nil
}

// PushFunc 推入一条函数式任务
func (d *dispatcher) PushFunc(f func() error) error {
	if !d.running {
		return nil
	}

	d.jobQueue <- JobFunc(f)
	return nil
}

func (d *dispatcher) stop() {
	for i := 0; i < len(d.workers); i++ {
		d.workers[i].Stop()
	}

	close(d.workPool)

	close(d.jobQueue)

	log.Println("dispatcher.stop()")
}

// Stop 停止任务
func (d *dispatcher) Stop() error {
	d.running = false
	d.quit <- true
	return nil
}

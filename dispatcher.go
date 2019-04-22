package gowork

import (
	"errors"
	"sync"
)

// Dispatcher 分配器
type Dispatcher interface {
	Push(Job) error
	PushFunc(func() error) error

	Run() error
	Stop() error
	Running() bool
}

// NewDispatcher 返回一个新的分配器
func NewDispatcher(cap int) Dispatcher {
	if cap <= 0 {
		panic(errors.New("cap can't not less than zero"))
	}

	d := &dispatcher{
		cap: cap,

		workers:  make([]Worker, 0, cap),

		jobQueue: make(jobChan),

		workPool: make(workerChan, cap),

		quit: make(chan bool),

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

	runLocker sync.Mutex
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
	d.runLocker.Lock()
	defer d.runLocker.Unlock()

	if d.running {
		return nil
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

	d.running = true

	d.run()

	return nil
}

func (d *dispatcher) stop() {
	for i := 0; i < len(d.workers); i++ {
		d.workers[i].Stop()
	}
}

// Stop 停止任务
func (d *dispatcher) Stop() error {
	if !d.running {
		return nil
	}
	d.running = false
	d.quit <- true
	return nil
}

// Running 任务是否正在执行
func (d *dispatcher) Running() bool{
	return d.running
}

// Push 推入一条任务
func (d *dispatcher) Push(j Job) error {
	if !d.running {
		return errors.New("dispatcher is no running")
	}
	if len(d.workers) < d.cap && len(d.workPool) == 0 {
		w := newWorker(d)
		d.workers = append(d.workers, w)
		if d.running {
			w.Start()
		}
	}

	go func() {
		d.jobQueue <- j
	}()
	return nil
}

// PushFunc 推入一条函数式任务
func (d *dispatcher) PushFunc(f func() error) error {
	return d.Push(JobFunc(f))
}

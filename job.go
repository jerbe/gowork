package go_pool

// Job 任务
type Job interface {
	Execute() error
}

type jobChan chan Job

// JobFunc 任务方法
type JobFunc func() error

// Execute 执行任务
func (j JobFunc) Execute() error {
	return j()
}

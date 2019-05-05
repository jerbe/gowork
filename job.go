package gowork

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

var (
	// 执行错误时调用的方法
	ExecuteErrorHandle func(error)
)

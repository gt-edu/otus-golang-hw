package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type RunResult struct {
	executedWorkersCount int32
	errCount             int32
	err                  error
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) RunResult {
	errBufCh := make(chan int32, 1)
	completeCh := make(chan struct{})
	tasksBufCh := make(chan Task, len(tasks))

	res := RunResult{}

	wg := sync.WaitGroup{}
	for i := 0; i < n; i++ {
		wg.Add(1)
		go worker(&wg, tasksBufCh, m, errBufCh, &res)
	}

	go func() {
		wg.Wait()
		completeCh <- struct{}{}
		close(completeCh)
	}()

	for _, t := range tasks {
		tasksBufCh <- t
	}
	close(tasksBufCh)

	<-completeCh

	var retErr error
	select {
	case <-errBufCh:
		retErr = ErrErrorsLimitExceeded
	default:
	}
	res.err = retErr
	return res
}

func worker(wg *sync.WaitGroup, tasksCh chan Task, maxErrors int, errBufCh chan int32, result *RunResult) {
	defer wg.Done()

	once := sync.Once{}
	for {
		task, ok := <-tasksCh
		if !ok {
			return
		}
		once.Do(func() {
			atomic.AddInt32(&result.executedWorkersCount, 1)
		})

		if maxErrors > 0 && atomic.LoadInt32(&result.errCount) >= int32(maxErrors) {
			return
		}

		err := task()
		if maxErrors > 0 && err != nil {
			if currErrCount := atomic.AddInt32(&result.errCount, 1); currErrCount == int32(maxErrors) {
				errBufCh <- currErrCount
				close(errBufCh)
			}
		}
	}
}

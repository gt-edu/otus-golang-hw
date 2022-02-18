package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	errBufCh := make(chan int32, 1)
	completeCh := make(chan struct{})
	tasksCount := len(tasks)
	tasksBufCh := make(chan Task, tasksCount)

	var wg sync.WaitGroup
	var errCount int32

	workersCount := n
	if workersCount > tasksCount {
		workersCount = tasksCount
	}
	for i := 0; i < n; i++ {
		wg.Add(1)
		go worker(&wg, tasksBufCh, m, errBufCh, &errCount)
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
	return retErr
}

func worker(wg *sync.WaitGroup, tasksCh chan Task, maxErrors int, errBufCh chan int32, ptrErrCount *int32) {
	defer wg.Done()

	for {
		task, ok := <-tasksCh
		if !ok {
			return
		}

		if maxErrors > 0 && atomic.LoadInt32(ptrErrCount) >= int32(maxErrors) {
			return
		}

		err := task()
		if maxErrors > 0 && err != nil {
			if currErrCount := atomic.AddInt32(ptrErrCount, 1); currErrCount == int32(maxErrors) {
				errBufCh <- currErrCount
				close(errBufCh)
			}
		}
	}
}

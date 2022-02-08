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
	var errCount int32
	errChBuffered := make(chan int32, 1)
	completeCh := make(chan struct{})
	tasksChBuffered := make(chan Task, len(tasks))

	wg := sync.WaitGroup{}
	for i := 0; i <= n; i++ {
		wg.Add(1)
		go worker(&wg, tasksChBuffered, &errCount, m, errChBuffered)
	}

	go func() {
		wg.Wait()
		completeCh <- struct{}{}
		close(completeCh)
	}()

	for _, t := range tasks {
		tasksChBuffered <- t
	}
	close(tasksChBuffered)

	<-completeCh

	var retErr error
	select {
	case <-errChBuffered:
		retErr = ErrErrorsLimitExceeded
	default:
	}

	return retErr
}

func worker(wg *sync.WaitGroup, tasksCh chan Task, errCount *int32, m int, errChBuffered chan int32) {
	defer wg.Done()

	for {
		task, ok := <-tasksCh

		if !ok || atomic.LoadInt32(errCount) >= int32(m) {
			return
		}

		err := task()
		if err != nil {
			if currErr := atomic.AddInt32(errCount, 1); currErr == int32(m) {
				errChBuffered <- currErr
				close(errChBuffered)
			}
		}
	}
}

package hw06pipelineexecution

import (
	"sync"
	"sync/atomic"
)

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	var doneReceived int32
	newData := make([]interface{}, 0)
	newDataMutex := sync.Mutex{}

	var wg sync.WaitGroup
	i := 0
	for val := range in {
		newDataMutex.Lock()
		newData = append(newData, nil)
		newDataMutex.Unlock()

		// Не уверен, что верное решение, что нужно пересоздавать канал для одного элемента
		singleValueCh := make(Bi)

		wg.Add(1)
		go func(singleValueCh In, num int, wg *sync.WaitGroup) {
			defer wg.Done()

			for _, s := range stages {
				select {
				case <-done:
					atomic.AddInt32(&doneReceived, 1)
					return
				case v := <-singleValueCh:
					// Не уверен, что верное решение, что нужно пересоздавать канал для одного элемента
					newSingleValueCh := make(Bi)
					singleValueCh = s(newSingleValueCh)

					newSingleValueCh <- v
					close(newSingleValueCh)
				}
			}

			res := <-singleValueCh

			newDataMutex.Lock()
			newData[num] = res
			newDataMutex.Unlock()
		}(singleValueCh, i, &wg)
		i++

		singleValueCh <- val
		close(singleValueCh)
	}

	wg.Wait()

	out := make(Bi)

	if doneReceived == 0 {
		go func() {
			// Не уверен в оптимальности этого решения, но пока не нашел способа получить данные в той же
			// последовательности
			for _, v := range newData {
				out <- v
			}
			close(out)
		}()
	} else {
		close(out)
	}

	return out
}

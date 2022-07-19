package hw06pipelineexecution

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	sleepPerStage = time.Millisecond * 100
	fault         = sleepPerStage / 2
)

func TestPipeline(t *testing.T) {
	// Stage generator
	g := func(_ string, f func(v interface{}) interface{}) Stage {
		return func(in In) Out {
			out := make(Bi)
			go func() {
				defer close(out)
				for v := range in {
					time.Sleep(sleepPerStage)
					fmt.Printf("Run stage func: %s %v\n", time.Now().Format("15:04:05.000000"), v)
					out <- f(v)
				}
			}()
			return out
		}
	}

	stages := []Stage{
		g("Dummy", func(v interface{}) interface{} { return v }),
		g("Multiplier (* 2)", func(v interface{}) interface{} { return v.(int) * 2 }),
		g("Adder (+ 100)", func(v interface{}) interface{} { return v.(int) + 100 }),
		g("Stringifier", func(v interface{}) interface{} { return strconv.Itoa(v.(int)) }),
	}

	t.Run("simple case", func(t *testing.T) {
		in := make(Bi)
		data := []int{1, 3, 5, 7, 9}

		go func() {
			for _, v := range data {
				in <- v
			}
			close(in)
		}()

		result := make([]string, 0, 10)
		start := time.Now()
		for s := range ExecutePipeline(in, nil, stages...) {
			fmt.Println("One final value: ", time.Now().Format("15:04:05.000000"))
			result = append(result, s.(string))
		}
		elapsed := time.Since(start)

		require.Equal(t, []string{"102", "106", "110", "114", "118"}, result)
		require.Less(t,
			int64(elapsed),
			// ~0.8s for processing 5 values in 4 stages (100ms every) concurrently
			int64(sleepPerStage)*int64(len(stages)+len(data)-1)+int64(fault))
	})

	//t.Run("done case", func(t *testing.T) {
	//	in := make(Bi)
	//	done := make(Bi)
	//	data := []int{1, 2, 3, 4, 5}
	//
	//	// Abort after 200ms
	//	abortDur := sleepPerStage * 2
	//	go func() {
	//		<-time.After(abortDur)
	//		close(done)
	//	}()
	//
	//	go func() {
	//		for _, v := range data {
	//			in <- v
	//		}
	//		close(in)
	//	}()
	//
	//	result := make([]string, 0, 10)
	//	start := time.Now()
	//	for s := range ExecutePipeline(in, done, stages...) {
	//		result = append(result, s.(string))
	//	}
	//	elapsed := time.Since(start)
	//
	//	require.Len(t, result, 0)
	//	require.Less(t, int64(elapsed), int64(abortDur)+int64(fault))
	//})
}

package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 23
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.NoError(t, err)

		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})
}

func TestRunWithMixedCountErrorsMoreThanM(t *testing.T) {
	tests := []struct {
		tasksCount     int
		workersCount   int
		maxErrorsCount int
	}{
		{tasksCount: 50, workersCount: 4, maxErrorsCount: 4},
		{tasksCount: 50, workersCount: 4, maxErrorsCount: 24},
		{tasksCount: 50, workersCount: 4, maxErrorsCount: 4},
		{tasksCount: 50, workersCount: 4, maxErrorsCount: 4},
	}

	for _, tt := range tests {
		ttp := tt
		t.Run(fmt.Sprintf("mixed tasks %d %d %d", ttp.tasksCount, ttp.workersCount, ttp.maxErrorsCount), func(t *testing.T) {
			t.Parallel()
			tasksCount := ttp.tasksCount
			workersCount := ttp.workersCount
			maxErrorsCount := ttp.maxErrorsCount

			tasks, ptrRunTasksCount := makeTestTasksHalfWithErrors(tasksCount)

			err := Run(tasks, workersCount, maxErrorsCount)

			require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
			require.LessOrEqual(t, int32(maxErrorsCount), *ptrRunTasksCount)
			require.LessOrEqual(t, *ptrRunTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
		})
	}
}

func TestRunZeroErrorsOrLimitIsHigh(t *testing.T) {
	tests := []struct {
		tasksCount     int
		workersCount   int
		maxErrorsCount int
	}{
		{tasksCount: 50, workersCount: 1, maxErrorsCount: 0},
		{tasksCount: 50, workersCount: 4, maxErrorsCount: 0},
		{tasksCount: 50, workersCount: 50, maxErrorsCount: 0},
		{tasksCount: 50, workersCount: 60, maxErrorsCount: 0},
		{tasksCount: 50, workersCount: 4, maxErrorsCount: 0},
		{tasksCount: 50, workersCount: 5, maxErrorsCount: 0},
		{tasksCount: 10, workersCount: 20, maxErrorsCount: 0},
		{tasksCount: 10, workersCount: 10, maxErrorsCount: 0},
		{tasksCount: 10, workersCount: 9, maxErrorsCount: 0},
		{tasksCount: 50, workersCount: 4, maxErrorsCount: 26},
		{tasksCount: 50, workersCount: 4, maxErrorsCount: 30},
		{tasksCount: 50, workersCount: 4, maxErrorsCount: 35},
		{tasksCount: 50, workersCount: 4, maxErrorsCount: 50},
		{tasksCount: 50, workersCount: 4, maxErrorsCount: 60},
		{tasksCount: 50, workersCount: 5, maxErrorsCount: 60},
		{tasksCount: 10, workersCount: 20, maxErrorsCount: 60},
		{tasksCount: 10, workersCount: 10, maxErrorsCount: 60},
		{tasksCount: 10, workersCount: 9, maxErrorsCount: 60},
	}

	for _, tt := range tests {
		ttp := tt
		t.Run(fmt.Sprintf("mixed tasks %d %d %d", ttp.tasksCount, ttp.workersCount, ttp.maxErrorsCount), func(t *testing.T) {
			t.Parallel()
			tasksCount := ttp.tasksCount
			workersCount := ttp.workersCount
			maxErrorsCount := ttp.maxErrorsCount
			tasks, ptrRunTasksCount := makeTestTasksHalfWithErrors(tasksCount)

			err := Run(tasks, workersCount, maxErrorsCount)

			require.NoError(t, err)
			require.Equal(t, int32(tasksCount), *ptrRunTasksCount)
		})
	}
}

func TestRunTasksNoErrors(t *testing.T) {
	tests := []struct {
		tasksCount     int
		workersCount   int
		maxErrorsCount int
	}{
		{tasksCount: 50, workersCount: 5, maxErrorsCount: 1},
		{tasksCount: 10, workersCount: 20, maxErrorsCount: 1},
		{tasksCount: 10, workersCount: 10, maxErrorsCount: 1},
		{tasksCount: 10, workersCount: 9, maxErrorsCount: 1},
	}

	for _, tt := range tests {
		ttp := tt
		t.Run(fmt.Sprintf("mixed tasks %d %d %d", ttp.tasksCount, ttp.workersCount, ttp.maxErrorsCount), func(t *testing.T) {
			t.Parallel()
			tasksCount := ttp.tasksCount
			workersCount := ttp.workersCount
			maxErrorsCount := ttp.maxErrorsCount
			tasks, ptrRunTasksCount := makeTestTasksNoErrors(tasksCount)

			err := Run(tasks, workersCount, maxErrorsCount)

			require.NoError(t, err)
			require.Equal(t, int32(tasksCount), *ptrRunTasksCount)
		})
	}
}

func makeTestTasksHalfWithErrors(tasksCount int) ([]Task, *int32) {
	tasks := make([]Task, 0, tasksCount)
	var runTasksCount int32

	i := 0
	for ; i < tasksCount/2; i++ {
		err := fmt.Errorf("error from task %d", i)
		tasks = append(tasks, func() error {
			atomic.AddInt32(&runTasksCount, 1)
			return err
		})
	}
	for ; i < tasksCount; i++ {
		tasks = append(tasks, func() error {
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(20)))
			atomic.AddInt32(&runTasksCount, 1)
			return nil
		})
	}
	return tasks, &runTasksCount
}

func makeTestTasksNoErrors(tasksCount int) ([]Task, *int32) {
	tasks := make([]Task, 0, tasksCount)
	var runTasksCount int32

	i := 0
	for ; i < tasksCount; i++ {
		tasks = append(tasks, func() error {
			atomic.AddInt32(&runTasksCount, 1)
			return nil
		})
	}
	return tasks, &runTasksCount
}

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

	t.Run("no tasks", func(t *testing.T) {
		tasks := []Task{}
		workersCount := 5
		maxErrorsCount := 1

		err := Run(tasks, workersCount, maxErrorsCount)
		require.NoError(t, err, "empty task list should not return an error")
	})

	t.Run("single task with no errors", func(t *testing.T) {
		tasks := []Task{
			func() error {
				return nil
			},
		}
		workersCount := 5
		maxErrorsCount := 1

		err := Run(tasks, workersCount, maxErrorsCount)
		require.NoError(t, err, "single task should not return an error")
	})

	t.Run("concurrency check", func(t *testing.T) {
		var maxConcurrency int32
		var currentConcurrency int32

		tasks := make([]Task, 100)
		for i := range tasks {
			tasks[i] = func() error {
				atomic.AddInt32(&currentConcurrency, 1)
				defer atomic.AddInt32(&currentConcurrency, -1)

				for {
					oldMax := atomic.LoadInt32(&maxConcurrency)
					current := atomic.LoadInt32(&currentConcurrency)
					if current > oldMax && atomic.CompareAndSwapInt32(&maxConcurrency, oldMax, current) {
						break
					} else if current <= oldMax {
						break
					}
				}

				time.Sleep(10 * time.Millisecond)
				return nil
			}
		}

		workersCount := 10
		err := Run(tasks, workersCount, 1)
		require.NoError(t, err)
		require.Equal(t, int32(workersCount), maxConcurrency, "not all workers ran concurrently")
	})
}

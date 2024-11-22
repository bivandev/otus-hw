package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func Run(tasks []Task, n int, m int) error {
	if n <= 0 {
		n = 1
	}

	var (
		wg             sync.WaitGroup
		errCounter     int32
		stopProcessing int32
		taskCh         = make(chan Task)
	)

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range taskCh {
				if atomic.LoadInt32(&stopProcessing) > 0 {
					return
				}

				if err := task(); err != nil {
					if m <= 0 {
						atomic.StoreInt32(&stopProcessing, 1)
						return
					}

					if atomic.AddInt32(&errCounter, 1) >= int32(m) {
						atomic.StoreInt32(&stopProcessing, 1)
					}
				}
			}
		}()
	}

	for _, task := range tasks {
		if atomic.LoadInt32(&stopProcessing) > 0 {
			break
		}
		taskCh <- task
	}

	close(taskCh)
	wg.Wait()

	if m > 0 && atomic.LoadInt32(&errCounter) >= int32(m) {
		return ErrErrorsLimitExceeded
	}

	if m <= 0 && atomic.LoadInt32(&errCounter) > 0 {
		return ErrErrorsLimitExceeded
	}

	return nil
}

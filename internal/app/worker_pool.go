package app

import (
	"sync"

	"github.com/cheggaaa/pb/v3"
)

type Task interface {
	Process()
}

type WorkerPool struct {
	Tasks           []Task
	workerCount     int
	taskChan        chan Task
	wg              sync.WaitGroup
}

func NewWorkerPool(bar *pb.ProgressBar, tasks []Task, provider *serviceProvider, workerCount int) *WorkerPool {
	return &WorkerPool{
		Tasks:           tasks,
		workerCount:     workerCount,
	}
}

func (wp *WorkerPool) worker() {
	for task := range wp.taskChan {
		task.Process()
	}
	wp.wg.Done()
}

func (wp *WorkerPool) Run() {
	wp.taskChan = make(chan Task, len(wp.Tasks))

	wp.wg.Add(wp.workerCount)

	for i := 0; i < wp.workerCount; i++ {
		go wp.worker()
	}

	for _, task := range wp.Tasks {
		wp.taskChan <- task
	}

	close(wp.taskChan)

	wp.wg.Wait()
}

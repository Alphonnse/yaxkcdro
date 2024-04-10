package app

import (
	"sync"

	"github.com/cheggaaa/pb/v3"
)

type Task interface {
	Process(serviceProvider *serviceProvider, bar *pb.ProgressBar)
}

type WorkerPool struct {
	Tasks           []Task
	serviceProvider *serviceProvider
	bar             *pb.ProgressBar
	workerCount     int
	taskChan        chan Task
	wg              sync.WaitGroup
}

func NewWorkerPool(bar *pb.ProgressBar, tasks []Task, provider *serviceProvider, workerCount int) *WorkerPool {
	return &WorkerPool{
		Tasks:           tasks,
		bar:             bar,
		serviceProvider: provider,
		workerCount:     workerCount,
	}
}

func (wp *WorkerPool) worker() {
	for task := range wp.taskChan {
		task.Process(wp.serviceProvider, wp.bar)
		wp.wg.Done()
	}
}

func (wp *WorkerPool) Run() {
	wp.taskChan = make(chan Task, len(wp.Tasks))

	for i := 0; i < wp.workerCount; i++ {
		go wp.worker()
	}

	wp.wg.Add(len(wp.Tasks))
	for _, task := range wp.Tasks {
		wp.taskChan <- task
	}

	close(wp.taskChan)

	wp.wg.Wait()
}

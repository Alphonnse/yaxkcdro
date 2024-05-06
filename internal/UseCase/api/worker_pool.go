package api

import (
	"sync"

	"github.com/cheggaaa/pb/v3"
)

type Task interface {
	process()
}

type workerPool struct {
	Tasks           []Task
	workerCount     int
	taskChan        chan Task
	wg              sync.WaitGroup
}

func NewWorkerPool(bar *pb.ProgressBar, tasks []Task, workerCount int) *workerPool {
	return &workerPool{
		Tasks:           tasks,
		workerCount:     workerCount,
	}
}

func (wp *workerPool) worker() {
	for task := range wp.taskChan {
		task.process()
	}
	wp.wg.Done()
}

func (wp *workerPool) Run() {
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

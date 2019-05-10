package async_task

import (
	"context"
	"time"

	"github.com/pkg/errors"
)

type (
	// Tasks in queue could distribute to workers.
	dispatcher struct {
		taskQueue    chan Task
		workerPool   chan chan Task
		workers      []*worker
		errorHandler ErrorHandler
		quit         chan bool
	}

	// A worker could directly process a task.
	worker struct {
		id           int
		workerPool   chan chan Task
		taskChannel  chan Task
		errorHandler ErrorHandler
		timeout      time.Duration
		quit         chan bool
	}
)

// Start a dispatcher
func (d *dispatcher) start() {
	for _, worker := range d.workers {
		go worker.start()
	}
	go d.dispatcher()
}

// Stop a dispatcher
func (d *dispatcher) stop() {
	for _, worker := range d.workers {
		worker.stop()
	}
	d.quit <- true
}

// Tasks will send to workers
func (d *dispatcher) dispatcher() {
	defer func() {
		if r := recover(); r != nil {
			d.errorHandler(errors.Wrap(r.(error), "[err] dispatcher panic"))
			go d.dispatcher()
		}
	}()

	for {
		select {
		case task := <-d.taskQueue:
			workerTaskChannel := <-d.workerPool
			workerTaskChannel <- task
		case <-d.quit:
			// delete all workers
			for len(d.workerPool) > 0 {
				<-d.workerPool
			}
			return
		}
	}
}

// Start a worker
func (w *worker) start() {
	defer func() {
		if r := recover(); r != nil {
			w.errorHandler(errors.Wrap(r.(error), "[err] worker panic"))
			go w.start()
		}
	}()

	// worker ready
	w.workerPool <- w.taskChannel
	for {
		select {
		case task := <-w.taskChannel:
			ctx, cancel := context.WithTimeout(context.Background(), w.timeout)
			if err := task.Process(ctx); err != nil {
				w.errorHandler(err)
			}
			cancel()
			w.workerPool <- w.taskChannel
		case <-w.quit: // end
			return
		}
	}
}

// Stop a worker
func (w *worker) stop() {
	w.quit <- true
}

package async_task

import (
	"context"
	"fmt"
	"time"
)

type (
	// It is executed if error is raised
	ErrorHandler func(error)

	// It is the task to dealt with
	Task interface {
		Process(ctx context.Context) error
	}

	// External interface
	Keeper interface {
		AddTask(ctx context.Context, task Task) error
	}

	// Internal interface
	keeper struct {
		errorHandler ErrorHandler
		queueSize    int
		workerSize   int
		timeout      time.Duration
		dispatcher   *dispatcher
	}
)

// Create a object of the Keeper
func NewAsyncTask(opts ...Option) (Keeper, error) {
	k := &keeper{}
	o := []Option{
		WithErrorHandlerOption(func(err error) {
			fmt.Printf("%+v \n", err)
		}),
		WithQueueSizeOption(1000),
		WithWorkerSizeOption(5),
		WithTimeoutOption(time.Duration(60 * time.Second)),
	}
	o = append(o, opts...)
	for _, opt := range o {
		opt.apply(k)
	}

	d, err := k.newDispatcher()
	if err != nil {
		return nil, err
	}
	k.dispatcher = d
	k.dispatcher.start()
	return Keeper(k), nil
}

// Add a task in asynchronously
func (k *keeper) AddTask(ctx context.Context, task Task) error {
	if ctx == nil || task == nil {
		return fmt.Errorf("[err] AddTask empty params")
	}

	// check context timeout
	select {
	case k.dispatcher.taskQueue <- task:
	case <-ctx.Done():
		return fmt.Errorf("[err] AddTask timeout")
	}
	return nil
}

// Create a dispatcher object
func (k *keeper) newDispatcher() (*dispatcher, error) {
	workerPool := make(chan chan Task, k.workerSize)
	var workers []*worker

	for i := 0; i < k.workerSize; i++ {
		worker := &worker{
			id:           i,
			workerPool:   workerPool,
			taskChannel:  make(chan Task),
			quit:         make(chan bool),
			errorHandler: k.errorHandler,
			timeout:      k.timeout,
		}
		workers = append(workers, worker)
	}

	return &dispatcher{
		taskQueue:    make(chan Task, k.queueSize),
		workerPool:   workerPool,
		workers:      workers,
		quit:         make(chan bool),
		errorHandler: k.errorHandler,
	}, nil
}

package bigquery

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/api/option"

	"cloud.google.com/go/bigquery"
)

type (

	// message
	Job struct {
		Msgs []*Message
	}

	Message struct {
		DatasetId string
		tableId   string
		Data      bigquery.ValueSaver
	}

	// dispatcher
	WorkerDispatcher struct {
		jobQueue   chan Job
		workerPool chan chan Job
		workers    []*Worker
		quit       chan bool
		errFunc    ErrorHandler
	}

	// worker
	Worker struct {
		client     *bigquery.Client
		workerPool chan chan Job
		jobChannel chan Job
		quit       chan bool
		errFunc    ErrorHandler
	}
)

func (d *WorkerDispatcher) addQueue(ctx context.Context, msgs []*Message) error {
	if ctx == nil || len(msgs) == 0 {
		return fmt.Errorf("[err] AddQueue empty params")
	}

	// check context timeout
	select {
	case d.jobQueue <- Job{Msgs: msgs}:
	case <-ctx.Done():
		return fmt.Errorf("[err] AddQueue timeout")
	}
	return nil
}

func (wd *WorkerDispatcher) start() {
	for _, worker := range wd.workers {
		go worker.start()
	}
	go wd.dispatcher()
}

func (wd *WorkerDispatcher) stop() {
	wd.quit <- true
	for _, worker := range wd.workers {
		worker.stop()
	}
}

func (wd *WorkerDispatcher) dispatcher() {
	defer func() {
		if r := recover(); r != nil {
			wd.errFunc(errors.Wrap(r.(error), "[err] dispatcher panic"))
			go wd.dispatcher()
		}
	}()
	for {
		select {
		case job := <-wd.jobQueue:
			go func(job Job) {
				workerJobChannel := <-wd.workerPool
				workerJobChannel <- job
			}(job)
		case <-wd.quit:
			return
		}
	}
}

func (w *Worker) start() {
	defer func() {
		if r := recover(); r != nil {
			w.errFunc(errors.Wrap(r.(error), "[err] worker panic"))
			go w.start()
		}
	}()

	for {
		w.workerPool <- w.jobChannel
		select {
		case job := <-w.jobChannel:
			for _, msg := range job.Msgs {
				if err := w.insert(msg); err != nil {
					w.errFunc(err)
				}
			}
		case <-w.quit:
			return
		}
	}
}

func (w *Worker) stop() {
	w.quit <- true
}

func (w *Worker) insert(msg *Message) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := w.client.Dataset(msg.DatasetId).Table(
		msg.tableId).Inserter().Put(ctx, msg.Data); err != nil {
		return errors.Wrap(err, "[err] insert")
	}
	return nil
}

func newWorker(cfg *Config, fn ErrorHandler, pool chan chan Job) (*Worker, error) {
	client, err := bigquery.NewClient(context.Background(),
		cfg.projectId,
		option.WithTokenSource(cfg.jwt.TokenSource(context.Background())))
	if err != nil {
		return nil, errors.Wrap(err, "[err]  newWorker fail client")
	}
	return &Worker{
		workerPool: pool,
		jobChannel: make(chan Job),
		client:     client,
		errFunc:    fn,
		quit:       make(chan bool),
	}, nil
}

func newWorkerDispatcher(cfg *Config, fn ErrorHandler, workerCount int) (*WorkerDispatcher, error) {
	if cfg == nil || fn == nil || workerCount == 0 {
		return nil, fmt.Errorf("[err] newWorkerDispatcher empty params")
	}

	workerPool := make(chan chan Job, workerCount)
	var workers []*Worker
	for i := 0; i < workerCount; i++ {
		worker, err := newWorker(cfg, fn, workerPool)
		if err != nil {
			return nil, errors.Wrap(err, "[err] newWorkerDispatcher newWorkerDispatcher fail")
		}
		workers = append(workers, worker)
	}

	return &WorkerDispatcher{
		jobQueue:   make(chan Job, 1000), // max queue 1000
		workerPool: workerPool,
		workers:    workers,
		quit:       make(chan bool),
		errFunc:    fn,
	}, nil
}

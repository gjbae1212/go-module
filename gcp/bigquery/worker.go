package bigquery

import (
	"context"
	"fmt"
	"net"
	"os"
	"syscall"
	"time"

	"github.com/gjbae1212/go-module/util"
	"github.com/pkg/errors"
	"google.golang.org/api/option"

	"cloud.google.com/go/bigquery"
)

type (

	// message
	Job struct {
		Msg *Message
	}

	Message struct {
		DatasetId string
		TableId   string
		Data      Row
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
		id         int
		client     *bigquery.Client
		workerPool chan chan Job
		jobChannel chan Job
		jobs       []Job
		maxStack   int
		delay      time.Duration
		quit       chan bool
		errFunc    ErrorHandler
	}
)

func (d *WorkerDispatcher) addQueue(ctx context.Context, msg *Message) error {
	if ctx == nil || msg == nil {
		return fmt.Errorf("[err] AddQueue empty params")
	}

	// check context timeout
	select {
	case d.jobQueue <- Job{Msg: msg}:
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
	for _, worker := range wd.workers {
		worker.stop()
	}
	wd.quit <- true
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
			workerJobChannel := <-wd.workerPool
			workerJobChannel <- job
		case <-wd.quit:
			// worker 다 제거
			for len(wd.workerPool) > 0 {
				<-wd.workerPool
			}
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

	// worker ready
	w.workerPool <- w.jobChannel
	for {
		select {
		case job := <-w.jobChannel: // job channel 에 job 이 들어 왔을때
			w.enqueue(job)
			if len(w.jobs) < w.maxStack {
				// worker ready
				w.workerPool <- w.jobChannel
				continue
			}

			// insert
			if errs := w.insertAll(); len(errs) > 0 {
				for _, err := range errs {
					w.errFunc(err)
				}
			}

			// worker ready
			w.workerPool <- w.jobChannel
		case <-time.After(w.delay): // delay 기간 동안 대기
			if errs := w.insertAll(); len(errs) > 0 {
				for _, err := range errs {
					w.errFunc(err)
				}
			}
		case <-w.quit: // 종료시
			if errs := w.insertAll(); len(errs) > 0 {
				for _, err := range errs {
					w.errFunc(err)
				}
			}
			return
		}
	}
}

func (w *Worker) stop() {
	w.quit <- true
}

func (w *Worker) enqueue(job ...Job) {
	w.jobs = append(w.jobs, job...)
}

func (w *Worker) insertAll() []error {
	var errs []error
	if len(w.jobs) == 0 {
		return errs
	}

	// wait max 1 minute
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	// retries logic
	var retries []Job
	// count
	total := len(w.jobs)
	fail := 0
	defer func() {
		w.jobs = w.jobs[:0]
		if len(retries) > 0 {
			fmt.Printf("[bq-worker-%d][retry] job retry! %d \n", w.id, len(retries))
			w.enqueue(retries...)
		}
	}()

	// job normalize
	categories := make(map[string]map[string][]Row)
	insertMap := make(map[string]Job)
	for _, job := range w.jobs {
		if _, ok := categories[job.Msg.DatasetId]; !ok {
			categories[job.Msg.DatasetId] = make(map[string][]Row)
		}
		categories[job.Msg.DatasetId][job.Msg.TableId] = append(categories[job.Msg.DatasetId][job.Msg.TableId], job.Msg.Data)
		if job.Msg.Data.InsertId() != "" {
			insertMap[job.Msg.Data.InsertId()] = job
		}
	}

	// insert all
	for datasetId, m := range categories {
		for tableId, rows := range m {
			inserter := w.client.Dataset(datasetId).Table(tableId).Inserter()
			inserter.SkipInvalidRows = true
			inserter.IgnoreUnknownValues = true
			if err := inserter.Put(ctx, rows); err != nil {
				if multierr, ok := err.(bigquery.PutMultiError); ok {
					for _, rowerr := range multierr {
						if len(rowerr.Errors) > 0 {
							fmt.Println(rowerr.Errors[0])
							errs = append(errs, rowerr.Errors[0])
							if w.isRetryable(rowerr.Errors[0]) {
								if job, ok := insertMap[rowerr.InsertID]; ok { // partly retry
									retries = append(retries, job)
								}
							} else {
								fail += 1
							}
						}
					}
				} else {
					fmt.Println(err)
					errs = append(errs, err)
					if w.isRetryable(err) {
						for _, row := range rows { // all retry
							retries = append(retries, Job{Msg: &Message{
								DatasetId: datasetId,
								TableId:   tableId,
								Data:      row,
							}})
						}
					} else {
						fail += len(rows)
					}
				}
			}
		}
	}

	fmt.Printf("[bq-worker-%d][%s] total %d insert %d fail %d retry %d \n", w.id, util.TimeToString(time.Now()),
		total, total-fail-len(retries), fail, len(retries))
	return errs
}

func (w *Worker) insert(msg *Message) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := w.client.Dataset(msg.DatasetId).Table(
		msg.TableId).Inserter().Put(ctx, msg.Data); err != nil {
		return errors.Wrap(err, "[err] insert")
	}
	return nil
}

func (w *Worker) isRetryable(err error) bool {
	switch err.(type) {
	case *bigquery.Error: // bigquery error timeout
		berr := err.(*bigquery.Error)
		if berr.Reason == "timeout" {
			return true
		}
	case *net.OpError: // catch connection reset
		neterr := err.(*net.OpError)
		// If this error message is `connection reset`, so jobs is retried.
		if syserr, ok := neterr.Err.(*os.SyscallError); ok && syserr.Err == syscall.ECONNRESET {
			return true
		}
	}
	if err == context.DeadlineExceeded || err == context.Canceled {
		return true
	}

	return false
}

func newWorker(id int, cfg *Config, fn ErrorHandler, pool chan chan Job) (*Worker, error) {
	client, err := bigquery.NewClient(context.Background(),
		cfg.projectId,
		option.WithTokenSource(cfg.jwt.TokenSource(context.Background())))
	if err != nil {
		return nil, errors.Wrap(err, "[err]  newWorker fail client")
	}
	return &Worker{
		id:         id,
		workerPool: pool,
		jobChannel: make(chan Job),
		jobs:       []Job{},
		maxStack:   cfg.workerStack,
		delay:      cfg.workerDelay,
		client:     client,
		errFunc:    fn,
		quit:       make(chan bool),
	}, nil
}

func newWorkerDispatcher(cfg *Config, fn ErrorHandler) (*WorkerDispatcher, error) {
	if cfg == nil || fn == nil {
		return nil, fmt.Errorf("[err] newWorkerDispatcher empty params")
	}

	workerPool := make(chan chan Job, cfg.workerSize)
	var workers []*Worker
	for i := 0; i < cfg.workerSize; i++ {
		worker, err := newWorker(i, cfg, fn, workerPool)
		if err != nil {
			return nil, errors.Wrap(err, "[err] newWorkerDispatcher newWorkerDispatcher fail")
		}
		workers = append(workers, worker)
	}

	return &WorkerDispatcher{
		jobQueue:   make(chan Job, cfg.queueSize),
		workerPool: workerPool,
		workers:    workers,
		quit:       make(chan bool),
		errFunc:    fn,
	}, nil
}

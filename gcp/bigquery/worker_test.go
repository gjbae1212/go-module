package bigquery

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"syscall"
	"testing"
	"time"

	"cloud.google.com/go/bigquery"

	"github.com/stretchr/testify/assert"
)

func TestNewWorkerDispatcher(t *testing.T) {
	assert := assert.New(t)

	cfg := testconfig()
	if cfg == nil {
		return
	}

	errFunc := func(err error) {
		log.Println(err)
	}

	dispatcher, err := newWorkerDispatcher(cfg, errFunc)
	assert.NoError(err)
	assert.Len(dispatcher.workers, 2)

	st := &streamer{}
	msg := &Message{
		DatasetId: cfg.schemas[0].DatasetId,
		TableId:   st.getTableId(cfg.schemas[0], time.Now().Add(24*time.Hour)),
		Data:      &TestItem{UserId: bigquery.NullInt64{Int64: 10}}}
	err = dispatcher.addQueue(context.Background(), msg)

	assert.NoError(err)
	assert.Len(dispatcher.jobQueue, 1)
	dispatcher.start()
	time.Sleep(2 * time.Second)
	assert.Len(dispatcher.jobQueue, 0)
	dispatcher.stop()

	err = dispatcher.addQueue(context.Background(), msg)
	time.Sleep(2 * time.Second)
	assert.Len(dispatcher.jobQueue, 1)
	dispatcher.start()
	time.Sleep(2 * time.Second)
	assert.Len(dispatcher.jobQueue, 0)

	// TODO: worker test
	for i := 0; i < 12321; i++ {
		msg := &Message{
			DatasetId: cfg.schemas[0].DatasetId,
			TableId:   st.getTableId(cfg.schemas[0], time.Now().Add(24*time.Hour)),
			Data:      &TestItem{UserId: bigquery.NullInt64{Int64: int64(i)}}}
		err = dispatcher.addQueue(context.Background(), msg)
	}
	time.Sleep(5 * time.Second)
}

func TestWorker(t *testing.T) {
	assert := assert.New(t)

	cfg := testconfig()
	if cfg == nil {
		return
	}

	errFunc := func(err error) {
		log.Println(err)
	}

	dispatcher, err := newWorkerDispatcher(cfg, errFunc)
	assert.NoError(err)
	worker := dispatcher.workers[0]

	item := &TestItem{UserId: bigquery.NullInt64{Int64: 9999999, Valid: true}}

	// enqueue
	worker.enqueue(Job{Msg: &Message{
		DatasetId: "empty",
		TableId:   "empty",
		Data:      item,
	}})
	assert.Len(worker.jobs, 1)

	// is retryable
	ok := worker.isRetryable(context.DeadlineExceeded)
	assert.True(ok)

	ok = worker.isRetryable(context.Canceled)
	assert.True(ok)

	merr := bigquery.PutMultiError{}
	ok = worker.isRetryable(merr)
	assert.False(ok)

	berr := &bigquery.Error{Reason: ""}
	ok = worker.isRetryable(berr)
	assert.False(ok)
	berr.Reason = "timeout"
	ok = worker.isRetryable(berr)
	assert.True(ok)

	nerr := &net.OpError{}
	nerr.Err = &os.SyscallError{Err: syscall.EACCES}
	ok = worker.isRetryable(nerr)
	assert.False(ok)
	nerr.Err = &os.SyscallError{Err: syscall.ECONNRESET}
	ok = worker.isRetryable(berr)
	assert.True(ok)
	ok = worker.isRetryable(fmt.Errorf("hello test error connection reset by peer hi"))
	assert.True(ok)

	// insert
	err = worker.insert(&Message{
		DatasetId: "empty",
		TableId:   "empty",
		Data:      item,
	})
	assert.Error(err)

	// insert all
	errs := worker.insertAll()
	assert.Len(errs, 1)
	assert.Len(worker.jobs, 0)

	st := &streamer{}
	msg := &Message{
		DatasetId: cfg.schemas[0].DatasetId,
		TableId:   st.getTableId(cfg.schemas[0], time.Now().Add(24*time.Hour)),
		Data:      &TestItem{UserId: bigquery.NullInt64{Int64: int64(20304050)}}}
	worker.enqueue(Job{Msg: msg})
	errs = worker.insertAll()
	assert.Len(errs, 0)
	assert.Len(worker.jobs, 0)
}

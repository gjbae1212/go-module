package bigquery

import (
	"context"
	"log"
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

	dispatcher, err := newWorkerDispatcher(cfg, errFunc, 2)
	assert.NoError(err)
	assert.Len(dispatcher.workers, 2)

	st := &streamer{}
	msgs := []*Message{&Message{
		DatasetId: cfg.schemas[0].DatasetId,
		TableId:   st.getTableId(cfg.schemas[0], time.Now().Add(24*time.Hour)),
		Data:      &TestItem{UserId: bigquery.NullInt64{Int64: 10}},
	}}

	err = dispatcher.addQueue(context.Background(), msgs)
	assert.NoError(err)
	assert.Len(dispatcher.jobQueue, 1)
	dispatcher.start()
	time.Sleep(2 * time.Second)
	assert.Len(dispatcher.jobQueue, 0)
	dispatcher.stop()

	err = dispatcher.addQueue(context.Background(), msgs)
	time.Sleep(2 * time.Second)
	assert.Len(dispatcher.jobQueue, 1)
	dispatcher.start()
	time.Sleep(2 * time.Second)
	assert.Len(dispatcher.jobQueue, 0)
}

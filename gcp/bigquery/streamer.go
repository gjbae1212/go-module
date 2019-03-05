package bigquery

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gjbae1212/go-module/util"

	"github.com/pkg/errors"

	"google.golang.org/api/option"

	"cloud.google.com/go/bigquery"
)

const (
	defaultQueueSize = 1000
)

type (
	Streamer interface {
		AddRow(ctx context.Context, row Row) error
	}

	streamer struct {
		cfg    *Config
		client *bigquery.Client

		async *WorkerDispatcher

		tickerLock sync.Mutex
		tickerDone chan *struct{}

		errFunc ErrorHandler
	}

	ErrorHandler func(error)
)

func NewStreamer(cfg *Config, errFunc ErrorHandler, queueSize int) (Streamer, error) {
	if cfg == nil {
		return nil, fmt.Errorf("[err] NewStreamerWithst empty params")
	}

	if errFunc == nil {
		errFunc = func(err error) {}
	}

	if queueSize < 0 {
		queueSize = defaultQueueSize
	}

	st := &streamer{cfg: cfg, errFunc: errFunc}

	// client 생성
	client, err := bigquery.NewClient(context.Background(),
		st.cfg.projectId,
		option.WithTokenSource(st.cfg.jwt.TokenSource(context.Background())))
	if err != nil {
		return nil, errors.Wrap(err, "[err]  NewStreamer fail client")
	}
	st.client = client

	dispatcher, err := newWorkerDispatcher(st.cfg, st.errFunc, 10, queueSize)
	if err != nil {
		return nil, errors.Wrap(err, "[err]  NewStreamer fail dispatcher")
	}
	st.async = dispatcher

	// start go routine
	st.async.start()
	st.ticker()
	return st, nil
}

func (st *streamer) AddRow(ctx context.Context, row Row) error {
	if row == nil || row.PublishedAt().IsZero() {
		return fmt.Errorf("[err] AddRow empty params")
	}

	schema, err := row.Schema()
	if err != nil {
		return errors.Wrap(err, "[err] AddRow unknown schema")
	}

	tableId := st.getTableId(schema, row.PublishedAt())
	var msgs []*Message
	msgs = append(msgs, &Message{
		DatasetId: schema.DatasetId,
		TableId:   tableId,
		Data:      row,
	})
	return st.async.addQueue(ctx, msgs)
}

func (st *streamer) ticker() error {
	st.tickerLock.Lock()
	defer st.tickerLock.Unlock()

	if st.tickerDone != nil {
		return nil
	}

	st.tickerDone = make(chan *struct{})
	go func(done chan *struct{}, errFunc ErrorHandler) {
		hourTicker := time.NewTicker(1 * time.Hour)
		if err := st.createTable(); err != nil {
			errFunc(err)
		}
		for {
			select {
			case <-done:
				return
			case <-hourTicker.C:
				if err := st.createTable(); err != nil {
					errFunc(err)
				}
			}
		}
	}(st.tickerDone, st.errFunc)
	return nil
}

func (st *streamer) deleteTicker() {
	st.tickerLock.Lock()
	defer st.tickerLock.Unlock()

	if st.tickerDone != nil {
		close(st.tickerDone)
		st.tickerDone = nil
	}
}

func (st *streamer) createTable() error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	for _, schema := range st.cfg.schemas {
		// 내일 테이블을 미리 생성
		tableId := st.getTableId(schema, time.Now().Add(24*time.Hour))
		table := st.client.Dataset(schema.DatasetId).Table(tableId)

		// table 없다면 error
		md, err := table.Metadata(ctx)
		if err != nil || md == nil {
			if err := table.Create(ctx,
				&bigquery.TableMetadata{Schema: schema.Schema}); err != nil {
				return errors.Wrap(err, "[err] createTable")
			}
		}
	}
	return nil
}

func (st *streamer) getTableId(schema *TableSchema, t time.Time) string {
	switch schema.Period {
	case Daily:
		return schema.Prefix + util.TimeToDailyStringFormat(t)
	case Monthly:
		return schema.Prefix + util.TimeToMonthlyStringFormat(t)
	}
	return ""
}

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

type (
	Streamer interface {
		AddRow(ctx context.Context, row Row) error
		AddRowSync(ctx context.Context, row Row) error
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

func NewStreamer(cfg *Config, errFunc ErrorHandler) (Streamer, error) {
	if cfg == nil {
		return nil, fmt.Errorf("[err] NewStreamerWithst empty params")
	}

	if errFunc == nil {
		errFunc = func(err error) {}
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

	dispatcher, err := newWorkerDispatcher(st.cfg, st.errFunc)
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

	return st.async.addQueue(ctx, &Message{
		DatasetId: schema.DatasetId,
		TableId:   st.getTableId(schema, row.PublishedAt()),
		Data:      row,
	})
}

// It is a function what data could insert into bigquery and waited until it is completed.
func (st *streamer) AddRowSync(ctx context.Context, row Row) error {
	if row == nil || row.PublishedAt().IsZero() {
		return fmt.Errorf("[err] AddRowSync empty params")
	}

	schema, err := row.Schema()
	if err != nil {
		return errors.Wrap(err, "[err] AddRowSync unknown schema")
	}

	inserter := st.client.Dataset(schema.DatasetId).Table(
		st.getTableId(schema, row.PublishedAt())).Inserter()
	inserter.SkipInvalidRows = true
	inserter.IgnoreUnknownValues = true
	return inserter.Put(ctx, row)
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
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// create today and tomorrow tables
	for _, schema := range st.cfg.schemas {
		for _, t := range []time.Time{time.Now(), time.Now().Add(24 * time.Hour)} {
			tableId := st.getTableId(schema, t)
			table := st.client.Dataset(schema.DatasetId).Table(tableId)
			md, err := table.Metadata(ctx)
			if err != nil || md == nil {
				if err := table.Create(ctx,
					&bigquery.TableMetadata{Schema: schema.Meta.Schema}); err != nil {
					return errors.Wrap(err, "[err] createTable")
				} else {
					fmt.Printf("[bq-table][%s] create %s\n", util.TimeToString(time.Now()), tableId)
				}
			}
		}
	}
	return nil
}

func (st *streamer) getTableId(schema *TableSchema, t time.Time) string {
	switch schema.Period {
	case NotExist:
		return schema.Prefix
	case Daily:
		return schema.Prefix + util.TimeToDailyStringFormat(t)
	case Monthly:
		return schema.Prefix + util.TimeToMonthlyStringFormat(t)
	case Yearly:
		return schema.Prefix + util.TimeToYearlyStringFormat(t)
	}
	return ""
}

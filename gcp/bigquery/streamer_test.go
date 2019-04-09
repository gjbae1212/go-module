package bigquery

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/option"

	"github.com/stretchr/testify/assert"
)

type TestItem struct {
	UserId bigquery.NullInt64
}

// Save implements the ValueSaver interface.
func (i *TestItem) Save() (map[string]bigquery.Value, string, error) {
	// unique value 넣여믄 중복 값 있으면 다시 넣지 않음.
	return map[string]bigquery.Value{
		"UserId": i.UserId.Int64,
		"Ignore": "ignore?",
	}, fmt.Sprintf("%d", i.UserId.Int64), nil
}

func (i *TestItem) Schema() (*TableSchema, error) {
	schema, err := bigquery.InferSchema(i)
	if err != nil {
		return nil, err
	}

	return &TableSchema{
		DatasetId: "test_dataset",
		Prefix:    "test_table_",
		Meta: &bigquery.TableMetadata{Schema: schema},
		Period: Daily,
	}, nil
}

func (i *TestItem) PublishedAt() time.Time {
	return time.Now().Add(time.Hour * 24)
}

type TestItem2 struct {
	UserId bigquery.NullInt64
}

// Save implements the ValueSaver interface.
func (i *TestItem2) Save() (map[string]bigquery.Value, string, error) {
	// unique value 넣여믄 중복 값 있으면 다시 넣지 않음.
	return map[string]bigquery.Value{
		"UserId": i.UserId.Int64,
		"Ignore": "ignore?",
	}, fmt.Sprintf("%d", i.UserId.Int64), nil
}

func (i *TestItem2) Schema() (*TableSchema, error) {
	schema, err := bigquery.InferSchema(i)
	if err != nil {
		return nil, err
	}

	return &TableSchema{
		DatasetId: "test_dataset",
		Prefix:    "allan_table",
		Meta: &bigquery.TableMetadata{Schema: schema},
		Period: NotExist,
	}, nil
}

func (i *TestItem2) PublishedAt() time.Time {
	return time.Now().Add(time.Hour * 24)
}

func TestTicker(t *testing.T) {
	assert := assert.New(t)

	cfg := testconfig()
	if cfg == nil {
		return
	}

	daily := &streamer{cfg: cfg, errFunc: func(err error) {
		log.Println(err)
	}}

	client, err := bigquery.NewClient(context.Background(),
		daily.cfg.projectId,
		option.WithTokenSource(daily.cfg.jwt.TokenSource(context.Background())))

	assert.NoError(err)
	daily.client = client

	daily.ticker()
	time.Sleep(3 * time.Second)
	daily.deleteTicker()
	time.Sleep(3 * time.Second)
}

func TestAsync(t *testing.T) {
	assert := assert.New(t)

	cfg := testconfig()
	if cfg == nil {
		return
	}

	daily := &streamer{cfg: cfg, errFunc: func(err error) {
		log.Println(err)
	}}

	client, err := bigquery.NewClient(context.Background(),
		daily.cfg.projectId,
		option.WithTokenSource(daily.cfg.jwt.TokenSource(context.Background())))

	assert.NoError(err)
	daily.client = client

	dispatcher, err := newWorkerDispatcher(daily.cfg, daily.errFunc)
	assert.NoError(err)
	daily.async = dispatcher
	daily.async.start()

	item := &TestItem{UserId: bigquery.NullInt64{Int64: 1}}
	err = daily.AddRow(context.Background(), item)
	assert.NoError(err)
	time.Sleep(5 * time.Second)
}

func TestStreamer_AddRow(t *testing.T) {
	assert := assert.New(t)

	cfg := testconfig()
	if cfg == nil {
		return
	}

	daily := &streamer{cfg: cfg, errFunc: func(err error) {
		log.Println(err)
	}}

	client, err := bigquery.NewClient(context.Background(),
		daily.cfg.projectId,
		option.WithTokenSource(daily.cfg.jwt.TokenSource(context.Background())))

	assert.NoError(err)
	daily.client = client

	dispatcher, err := newWorkerDispatcher(daily.cfg, daily.errFunc)
	assert.NoError(err)
	daily.async = dispatcher

	err = daily.AddRow(context.Background(), nil)
	assert.Error(err)

	item := &TestItem{UserId: bigquery.NullInt64{Int64: 1}}
	for i := 0; i < 1000; i++ {
		err = daily.AddRow(context.Background(), item)
		assert.NoError(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	err = daily.AddRow(ctx, item)
	assert.Error(err)
}

func TestStreamer_AddRowSync(t *testing.T) {
	assert := assert.New(t)

	cfg := testconfig()
	if cfg == nil {
		return
	}

	daily := &streamer{cfg: cfg, errFunc: func(err error) {
		log.Println(err)
	}}

	client, err := bigquery.NewClient(context.Background(),
		daily.cfg.projectId,
		option.WithTokenSource(daily.cfg.jwt.TokenSource(context.Background())))

	assert.NoError(err)
	daily.client = client

	dispatcher, err := newWorkerDispatcher(daily.cfg, daily.errFunc)
	assert.NoError(err)
	daily.async = dispatcher

	err = daily.AddRowSync(context.Background(), nil)
	assert.Error(err)

	item := &TestItem2{UserId: bigquery.NullInt64{Int64: 1}}
	for i := 0; i < 10; i++ {
		item.UserId.Int64 = int64(i)
		err = daily.AddRowSync(context.Background(), item)
		assert.NoError(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	cancel()
	err = daily.AddRowSync(ctx, item)
	assert.Error(err)
}

func TestStreamer_GetTableId(t *testing.T) {
	assert := assert.New(t)

	cfg := testconfig()
	if cfg == nil {
		return
	}

	daily := &streamer{cfg: cfg, errFunc: func(err error) {
		log.Println(err)
	}}

	schema := &TableSchema{
		Prefix: "aa",
		Period: Monthly,
	}
	now := time.Now()
	id := daily.getTableId(schema, now)
	assert.Equal(fmt.Sprintf("aa%d%02d", now.Year(), now.Month()), id)
	log.Println(id)

	schema = &TableSchema{
		Prefix: "bb",
		Period: Daily,
	}
	id = daily.getTableId(schema, now)
	assert.Equal(fmt.Sprintf("bb%d%02d%02d", now.Year(), now.Month(), now.Day()), id)
	log.Println(id)

	schema = &TableSchema{
		Prefix: "allan",
		Period: NotExist,
	}
	id = daily.getTableId(schema, now)
	assert.Equal(schema.Prefix,id)
	log.Println(id)

}

func testconfig() *Config {
	jwtpath := os.Getenv("GCP_JWT")
	_, err := os.Stat(jwtpath)
	if os.IsNotExist(err) {
		return nil
	}

	projectId := os.Getenv("PROJECT_ID")
	datasetId := "test_dataset"
	schema1, err := bigquery.InferSchema(TestItem{})
	if err != nil {
		return nil
	}
	ss1 := &TableSchema{
		DatasetId: datasetId,
		Prefix:    "test_table_",
		Meta: &bigquery.TableMetadata{Schema: schema1},
		Period: Daily,
	}

	schema2, err := bigquery.InferSchema(TestItem2{})
	if err != nil {
		return nil
	}
	ss2 := &TableSchema{
		DatasetId: datasetId,
		Prefix:    "allan_table",
		Meta: &bigquery.TableMetadata{Schema: schema2},
		Period: NotExist,
	}

	jwt, err := ioutil.ReadFile(jwtpath)
	if err != nil {
		return nil
	}

	cfg, err := NewConfig(projectId, jwt, []*TableSchema{ss1, ss2}, 1000, 2, 500, 1*time.Second)
	if err != nil {
		return nil
	}
	return cfg
}

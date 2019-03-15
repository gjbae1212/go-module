package bigquery

import (
	"io/ioutil"
	"os"
	"testing"

	"cloud.google.com/go/bigquery"

	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	assert := assert.New(t)

	_, err := NewConfig("", nil, nil, 0, 0, 0, time.Second)
	assert.Error(err)

	jwtpath := os.Getenv("GCP_JWT")
	_, err = os.Stat(jwtpath)
	if os.IsNotExist(err) {
		return
	}

	projectId := os.Getenv("PROJECT_ID")
	datasetId := "test_dataset"

	jwt, err := ioutil.ReadFile(jwtpath)
	assert.NoError(err)

	schema := &TableSchema{
		DatasetId: datasetId,
		Prefix:    "test_table_",
		Meta:      &bigquery.TableMetadata{Schema: bigquery.Schema{}},
	}

	c, err := NewConfig(projectId, jwt, []*TableSchema{schema}, 0, 0, 0, time.Second)
	assert.Error(err)

	schema.Meta.Schema = append(schema.Meta.Schema, &bigquery.FieldSchema{})
	c, err = NewConfig(projectId, jwt, []*TableSchema{schema}, 0, 0, 0, time.Second)
	assert.NoError(err)
	assert.Equal(defaultQueueSize, c.queueSize)
	assert.Equal(defaultWorkerSize, c.workerSize)
	assert.Equal(defaultWorkerStack, c.workerStack)
	assert.Equal(time.Second, c.workerDelay)
}

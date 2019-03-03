package bigquery

import (
	"io/ioutil"
	"os"
	"testing"

	"cloud.google.com/go/bigquery"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	assert := assert.New(t)

	_, err := NewConfig("", nil, nil)
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
		Schema:    bigquery.Schema{},
	}

	_, err = NewConfig(projectId, jwt, []*TableSchema{schema})
	assert.NoError(err)
}

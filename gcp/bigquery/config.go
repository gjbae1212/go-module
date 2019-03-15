package bigquery

import (
	"fmt"

	"time"

	"cloud.google.com/go/bigquery"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
)

type TablePeriod int

const (
	defaultQueueSize   = 1000
	defaultWorkerSize  = 10
	defaultWorkerStack = 500 // bigquery document recommend size
)

const (
	Daily TablePeriod = iota
	Monthly
	Yearly
)

type Config struct {
	projectId   string         // gcp project
	jwt         *jwt.Config    // gcp jwt config
	schemas     []*TableSchema // table schemas
	queueSize   int            // max queue
	workerSize  int            // worker count
	workerStack int            // worker stack size
	workerDelay time.Duration  // worker insert wait duration
}

type TableSchema struct {
	DatasetId string                  // bigquery datasetId
	Prefix    string                  // bigquery table prefix
	Meta      *bigquery.TableMetadata // bigquery table meta
	Period    TablePeriod             // TablePeriod
}

func NewConfig(projectId string, jwtbys []byte, schemas []*TableSchema, queueSize, workerSize, workerStack int, workerDelay time.Duration) (*Config, error) {
	if len(jwtbys) == 0 || projectId == "" || len(schemas) == 0 {
		return nil, fmt.Errorf("[err] NewConfig empty params")
	}
	for _, schema := range schemas {
		if len(schema.Meta.Schema) == 0 {
			return nil, fmt.Errorf("[err] NewConfig empty schema")
		}
	}

	if queueSize <= 0 {
		queueSize = defaultQueueSize
	}

	if workerSize <= 0 {
		workerSize = defaultWorkerSize
	}

	if workerStack <= 0 {
		workerStack = defaultWorkerStack
	}

	jwt, err := google.JWTConfigFromJSON(jwtbys, bigquery.Scope)
	if err != nil {
		return nil, err
	}

	return &Config{
		projectId: projectId,
		jwt:       jwt, schemas: schemas,
		queueSize:   queueSize,
		workerSize:  workerSize,
		workerStack: workerStack,
		workerDelay: workerDelay,
	}, nil
}

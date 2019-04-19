package bigquery

import (
	"time"

	"cloud.google.com/go/bigquery"
)

type Row interface {
	Save() (row map[string]bigquery.Value, insertID string, err error)
	Schema() (schema *TableSchema, err error)
	PublishedAt() time.Time
	InsertId() string
}

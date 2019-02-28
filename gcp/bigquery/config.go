package bigquery

import (
	"fmt"

	"cloud.google.com/go/bigquery"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
)

type TablePeriod int

const (
	Daily TablePeriod = iota
	Monthly
)

type Config struct {
	projectId string         // gcp project
	datasetId string         // biquery dataset id
	jwt       *jwt.Config    // gcp jwt config
	schemas   []*TableSchema // table schemas

}

type TableSchema struct {
	Prefix string          // bigquery table prefix
	Schema bigquery.Schema // bigquery table schema
	Period      TablePeriod     // TablePeriod
}

func NewConfig(projectId, datasetId string, jwtbys []byte, schemas []*TableSchema) (*Config, error) {
	if len(jwtbys) == 0 || projectId == "" || datasetId == "" {
		return nil, fmt.Errorf("[err] NewConfig empty params")
	}

	jwt, err := google.JWTConfigFromJSON(jwtbys, bigquery.Scope)
	if err != nil {
		return nil, err
	}

	return &Config{
		projectId: projectId,
		datasetId: datasetId,
		jwt:       jwt, schemas: schemas}, nil
}

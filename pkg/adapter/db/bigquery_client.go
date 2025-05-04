package db

import (
	"context"

	"cloud.google.com/go/bigquery"
	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/helper"
)

// NewBigqueryClient bigqueryのクライアントを作成する
func NewBigqueryClient() (*bigquery.Client, error) {
	ctx := context.Background()
	bigquery, err := bigquery.NewClient(ctx, helper.GetProjectId())
	if err != nil {
		return nil, err
	}
	return bigquery, nil

}

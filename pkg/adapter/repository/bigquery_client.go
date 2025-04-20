package repository

import (
	"cloud.google.com/go/bigquery"
	"context"
	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/helper"
)

func NewBigqueryClient() (*bigquery.Client, error) {
	ctx := context.Background()
	bigquery, err := bigquery.NewClient(ctx, helper.GetProjectId())
	if err != nil {
		return nil, err
	}
	return bigquery, nil

}

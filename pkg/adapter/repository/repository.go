package repository

import (
	"cloud.google.com/go/bigquery"
	"context"
	"errors"
	"fmt"
	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/domain/models"
	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/domain/repository"
	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/helper"
	"google.golang.org/api/iterator"
	"log"
)

const (
	transactionTable       = "Transactions"
	transactionReportTable = "TransactionReports"
)

type transactionRepository struct {
	client *bigquery.Client
}

func NewTransactionRepository(client *bigquery.Client) repository.TransactionRepository {
	return &transactionRepository{
		client: client,
	}
}

type transactionReportRepository struct {
	client *bigquery.Client
}

func NewMonthlySummaryRepository(client *bigquery.Client) repository.TransactionReportRepository {
	return &transactionReportRepository{
		client: client,
	}
}

func (r *transactionRepository) SaveAll(ctx context.Context, models []*models.Transaction) error {
	if err := r.client.Dataset(helper.GetDatasetId()).Table(transactionTable).Inserter().Put(ctx, models); err != nil {
		return err
	}
	return nil

}

func (tr *transactionReportRepository) SaveAll(ctx context.Context, models []*models.TransactionReport) error {
	if err := tr.client.Dataset(helper.GetDatasetId()).Table(transactionReportTable).Inserter().Put(ctx, models); err != nil {
		return err
	}
	return nil
}

func (tr *transactionReportRepository) SelectByYearMonth(ctx context.Context, yearMonth string) ([]*models.TransactionReport, error) {
	sql := fmt.Sprintf(
		"SELECT Id FROM `%s.%s` WHERE YearMonth = @yearMonth",
		helper.GetDatasetId(), transactionReportTable,
	)
	query := tr.client.Query(sql)
	// クエリパラメータを設定
	query.Parameters = []bigquery.QueryParameter{
		{
			Name:  "yearMonth",
			Value: yearMonth,
		},
	}

	it, err := query.Read(ctx)
	if err != nil {
		fmt.Println("Error reading query results:", err)
		return nil, err
	}
	var res []*models.TransactionReport
	for {
		var r models.TransactionReport
		err := it.Next(&r)
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate result: %v", err)
		}
		res = append(res, &r)
	}
	return res, nil
}

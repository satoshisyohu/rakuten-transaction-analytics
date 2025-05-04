package repository

import (
	"context"

	"cloud.google.com/go/bigquery"
	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/domain/models"
	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/domain/repository"
	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/helper"
)

const (
	transactionTable = "Transactions"
)

// transactionRepository 明細を保存するリポジトリ
type transactionRepository struct {
	client *bigquery.Client
}

// NewTransactionRepository コンストラクタ
func NewTransactionRepository(client *bigquery.Client) repository.TransactionRepository {
	return &transactionRepository{
		client: client,
	}
}

// SaveAll 明細を保存する
func (r *transactionRepository) SaveAll(ctx context.Context, models []*models.Transaction) error {
	return r.client.Dataset(helper.GetDatasetId()).Table(transactionTable).Inserter().Put(ctx, models)
}

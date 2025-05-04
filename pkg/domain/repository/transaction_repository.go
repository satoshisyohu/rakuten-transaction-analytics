package repository

import (
	"context"

	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/domain/models"
)

type TransactionRepository interface {
	SaveAll(ctx context.Context, models []*models.Transaction) error
}

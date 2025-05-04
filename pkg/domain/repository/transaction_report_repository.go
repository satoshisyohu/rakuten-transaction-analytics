package repository

import (
	"context"

	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/domain/models"
)

type TransactionReportRepository interface {
	SaveAll(ctx context.Context, models []*models.TransactionReport) error
	SelectByYearMonth(ctx context.Context, yearMonth string) ([]*models.TransactionReport, error)
}

//go:generate mockgen -source=transaction_report_repository.go -destination=../../../internal/mocks/domain/repository/transaction_report_repository.go -package=mocks

package repository

import (
	"context"

	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/domain/models"
)

type TransactionReportRepository interface {
	SaveAll(ctx context.Context, models []*models.TransactionReport) error
	SelectByYearMonth(ctx context.Context, yearMonth string) ([]*models.TransactionReport, error)
}

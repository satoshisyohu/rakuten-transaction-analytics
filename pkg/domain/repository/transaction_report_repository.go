//go:generate mockgen -source=transaction_report_repository.go -destination=../../../internal/mocks/domain/repository/transaction_report_repository.go -package=mocks

package repository

import (
	"context"

	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/domain/models"
)

// TransactionReportRepository repository interface
type TransactionReportRepository interface {
	SaveAll(context.Context, []*models.TransactionReport) error
	SelectByYearMonth(context.Context, string) ([]*models.TransactionReport, error)
}

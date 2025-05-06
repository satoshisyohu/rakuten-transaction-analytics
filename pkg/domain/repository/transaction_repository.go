//go:generate mockgen -source=transaction_repository.go -destination=../../../internal/mocks/domain/repository/transaction_repository.go -package=mocks

package repository

import (
	"context"

	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/domain/models"
)

// TransactionRepository repository interface
type TransactionRepository interface {
	SaveAll(ctx context.Context, models []*models.Transaction) error
}

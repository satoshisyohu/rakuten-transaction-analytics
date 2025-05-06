//go:generate mockgen -source=analytics_transaction_usecase.go -destination=../../internal/mocks/infrastructure/analytics_transaction_usecase.go -package=mocks

package infrastructure

import (
	"context"

	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/handler/dto"
)

// IAnalyticsTransactionUsecase interface
type IAnalyticsTransactionUsecase interface {
	// Run 実行する
	Run(context.Context, *dto.TransactionRequest) (*dto.TransactionResponse, error)
}

//go:generate mockgen -source=analytics_transaction_usecase.go -destination=../../internal/mocks/infrastructure/analytics_transaction_usecase.go -package=mocks

package infrastructure

import (
	"context"
	"mime/multipart"

	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/handler/dto"
)

type IAnalyticsTransactionUsecase interface {
	// Run 実行する
	Run(ctx context.Context, req dto.TransactionRequest, files []*multipart.FileHeader) (*dto.TransactionResponse, error)

	// Execute 改良中
	Execute(ctx context.Context, req dto.TransactionRequest) (*dto.TransactionResponse, error)
}

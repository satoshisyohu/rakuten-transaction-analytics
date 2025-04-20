package infrastructure

import (
	"context"
	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/handler/dto"
	"mime/multipart"
)

type IAnalyticsTransactionUsecase interface {
	// Run 実行する
	Run(ctx context.Context, req dto.TransactionRequest, files []*multipart.FileHeader) (*dto.TransactionResponse, error)
}

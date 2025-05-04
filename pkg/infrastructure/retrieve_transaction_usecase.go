package infrastructure

import (
	"context"
	"mime/multipart"

	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/handler/dto"
)

type IRetrieveTransactionUsecase interface {
	// Run 実行する
	Run(ctx context.Context, files []*multipart.FileHeader) (*dto.RetrieveTransactionResponses, error)
}

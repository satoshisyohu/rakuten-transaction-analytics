package rule

import (
	"context"

	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/domain/aggregate"
)

type ICalculateScoreRule interface {
	// CalculateScore ファイルを読み込む
	CalculateScore(ctx context.Context, transactionReportDto *aggregate.TransactionReportDto) (float64, error)
}

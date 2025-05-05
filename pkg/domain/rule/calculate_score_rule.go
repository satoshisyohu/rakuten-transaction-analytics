//go:generate mockgen -source=calculate_score_rule.go -destination=../../../internal/mocks/domain/rule/calculate_score_rule.go -package=mocks

package rule

import (
	"context"

	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/domain/aggregate"
)

type ICalculateScoreRule interface {
	// CalculateScore ファイルを読み込む
	CalculateScore(context.Context, *aggregate.TransactionReportDto) (float64, error)
}

package rule

import (
	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/domain/models"
)

type ICalculateScoreRule interface {
	// CalculateScore ファイルを読み込む
	CalculateScore(transactionReport *models.TransactionReport) (float64, error)
}

package impl

import (
	"errors"
	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/code"
	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/domain/models"
	"math"
)

// 減点係数
const penaltyFactor = 50.0

type CalculateScoreRule struct {
}

func NewCalculateScoreRule() *CalculateScoreRule {
	return &CalculateScoreRule{}
}

// CalculateScore ルールに基づいてスコアを計算する
func (c *CalculateScoreRule) CalculateScore(tr *models.TransactionReport) (float64, error) {

	scoreMapping := c.createScore(tr)
	if err := c.validateScoreMapping(scoreMapping); err != nil {
		return 0, err
	}

	// 初期スコアおよび減点係数の設定
	score := 100.0

	for _, s := range scoreMapping {
		// スコアを算出して、減点する
		score -= c.calculatePenaltyScore(s)
	}
	// 呼び出し元で使用しているので、スコアをセットする
	tr.Score = score

	return score, nil
}

func (c *CalculateScoreRule) calculatePenaltyScore(score *models.Score) float64 {
	return math.Abs(score.ActualRatio-score.IdealRatio) * penaltyFactor
}

func (c *CalculateScoreRule) createScore(tr *models.TransactionReport) []*models.Score {
	// todo ここもDBに持つのが理想的なので改善の余地あり
	return []*models.Score{
		// 食費
		{
			IdealRatio:  0.2,
			ActualRatio: float64(tr.FoodExpenses) / float64(tr.BaseAmounts),
			Category:    code.Food,
		},
		// 固定費
		{
			IdealRatio:  0.45,
			ActualRatio: float64(tr.FixedCosts) / float64(tr.BaseAmounts),
			Category:    code.Fixed,
		},
		// 変動費
		{
			IdealRatio:  0.05,
			ActualRatio: float64(tr.VariableCosts) / float64(tr.BaseAmounts),
			Category:    code.Variable,
		},
		// コンビニ費用
		{
			IdealRatio:  0.0,
			ActualRatio: float64(tr.WasteExpenses) / float64(tr.BaseAmounts),
			Category:    code.Waste,
		},
		// その他費用
		{
			IdealRatio:  0.1,
			ActualRatio: float64(tr.OtherExpenses) / float64(tr.BaseAmounts),
			Category:    code.Other,
		},
		// 貯金費用
		{
			IdealRatio:  0.2,
			ActualRatio: float64(tr.Savings) / float64(tr.BaseAmounts),
			Category:    code.Savings,
		},
	}
}

func (c *CalculateScoreRule) validateScoreMapping(scoreMapping []*models.Score) error {
	var sumScore float64
	for _, score := range scoreMapping {
		sumScore += score.IdealRatio
	}
	if sumScore != 1.0 {
		return errors.New("理想的な割合の合計が1.0ではありません")
	}
	return nil
}

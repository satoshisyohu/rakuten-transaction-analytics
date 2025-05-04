package impl

import (
	"context"
	"errors"
	"math"

	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/code"
	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/domain/aggregate"
	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/domain/models"
	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/domain/repository"
)

// 減点係数
const penaltyFactor = 50.0

type CalculateScoreRule struct {
	rrr repository.RatioRuleRepository
}

func NewCalculateScoreRule(rrr repository.RatioRuleRepository) *CalculateScoreRule {
	return &CalculateScoreRule{rrr: rrr}
}

// CalculateScore ルールに基づいてスコアを計算する
func (c *CalculateScoreRule) CalculateScore(ctx context.Context, tr *aggregate.TransactionReportDto) (float64, error) {

	// scoreを取得する
	ratioRules, err := c.rrr.SelectAll(ctx)
	if err != nil {
		return 0, err
	}
	// scoreをバリデーションする
	if err = c.validateScoreMapping(ratioRules); err != nil {
		return 0, err
	}
	// スコアを計算する
	mappedScore := c.calculateActualScore(tr, ratioRules)

	// 初期スコアおよび減点係数の設定
	score := 100.0

	for _, s := range mappedScore {
		// スコアを算出して、減点する
		score -= c.calculatePenaltyScore(s)
	}

	if score < 0 {
		return 0, errors.New("スコアが0未満になりました")
	}

	return score, nil
}

func (c *CalculateScoreRule) calculatePenaltyScore(score *models.Score) float64 {
	return math.Abs(score.ActualRatio-score.IdealRatio) * penaltyFactor
}

func (c *CalculateScoreRule) calculateActualScore(tr *aggregate.TransactionReportDto, ratioRules []*models.RatioRule) []*models.Score {
	// todo ここもDBに持つのが理想的なので改善の余地あり

	var scores []*models.Score
	for _, r := range ratioRules {
		var actualRatio float64
		switch r.Category {
		case code.Food.String():
			actualRatio = float64(tr.FoodExpenses) / float64(tr.BaseAmounts)
		case code.Other.String():
			actualRatio = float64(tr.OtherExpenses) / float64(tr.BaseAmounts)
		case code.Waste.String():
			actualRatio = float64(tr.WasteExpenses) / float64(tr.BaseAmounts)
		case code.Fixed.String():
			actualRatio = float64(tr.FixedCosts) / float64(tr.BaseAmounts)
		case code.Variable.String():
			actualRatio = float64(tr.VariableCosts) / float64(tr.BaseAmounts)
		case code.Savings.String():
			actualRatio = float64(tr.Savings) / float64(tr.BaseAmounts)
		}
		scores = append(scores, &models.Score{
			IdealRatio:  r.IdealRatio,
			ActualRatio: actualRatio,
		})
	}

	return scores
}

func (c *CalculateScoreRule) validateScoreMapping(scoreMapping []*models.RatioRule) error {
	var sumScore float64
	for _, score := range scoreMapping {
		sumScore += score.IdealRatio
	}
	if sumScore != 1.0 {
		return errors.New("理想的な割合の合計が1.0ではありません")
	}
	return nil
}

package impl

import (
	"context"
	"errors"
	"log/slog"
	"math"

	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/code"
	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/domain/aggregate"
	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/domain/models"
	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/domain/repository"
)

const penaltyFactor = 50.0 // 減点係数

// CalculateScoreRule スコア計算ルール
type CalculateScoreRule struct {
	rrr repository.RatioRuleRepository
}

// NewCalculateScoreRule CalculateScoreRuleのファクトリ関数
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
	if err = c.validateScore(ratioRules); err != nil {
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
	slog.Info("your score", "score", score)
	return score, nil
}

// calculateRatio baseAmountに対するカテゴリの支出の値を計算する
func (c *CalculateScoreRule) calculateRatio(expenses, baseAmount int64) float64 {
	return float64(expenses) / float64(baseAmount)
}

// calculatePenaltyScore 減点スコアを計算する
func (c *CalculateScoreRule) calculatePenaltyScore(score *models.Score) float64 {
	return math.Abs(score.ActualRatio-score.IdealRatio) * penaltyFactor
}

// calculateActualScore 実際のスコアを計算する
func (c *CalculateScoreRule) calculateActualScore(tr *aggregate.TransactionReportDto, ratioRules []*models.RatioRule) (scores []*models.Score) {

	for _, r := range ratioRules {
		var actualRatio float64
		switch r.Category {
		case code.Food.String():
			actualRatio = c.calculateRatio(tr.FoodExpenses, tr.BaseAmounts)
		case code.Other.String():
			actualRatio = c.calculateRatio(tr.OtherExpenses, tr.BaseAmounts)
		case code.Waste.String():
			actualRatio = c.calculateRatio(tr.WasteExpenses, tr.BaseAmounts)
		case code.Fixed.String():
			actualRatio = c.calculateRatio(tr.FixedCosts, tr.BaseAmounts)
		case code.Variable.String():
			actualRatio = c.calculateRatio(tr.VariableCosts, tr.BaseAmounts)
		case code.Savings.String():
			actualRatio = c.calculateRatio(tr.Savings, tr.BaseAmounts)

		}
		scores = append(scores, &models.Score{
			IdealRatio:  r.IdealRatio,
			ActualRatio: actualRatio,
		})
	}

	return scores
}

// validateScore dbから取得したスコアをバリデーションする
func (c *CalculateScoreRule) validateScore(scoreMapping []*models.RatioRule) error {
	var sumScore float64
	for _, score := range scoreMapping {
		sumScore += score.IdealRatio
	}
	if sumScore != 1.0 {
		return errors.New("理想的な割合の合計が1.0ではありません")
	}
	return nil
}

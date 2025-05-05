package impl

import (
	"context"
	"errors"
	"testing"

	mocks "github.com/satoshisyohu/rakuten-transaction-analytics/internal/mocks/domain/repository"
	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/domain/aggregate"
	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/domain/models"
	"github.com/zeebo/assert"
	"go.uber.org/mock/gomock"
)

var normalRatioModels = []*models.RatioRule{
	{
		Category:   "waste",
		IdealRatio: 0.0,
	},
	{
		Category:   "variable",
		IdealRatio: 0.05,
	},
	{
		Category:   "other",
		IdealRatio: 0.1,
	},
	{
		Category:   "savings",
		IdealRatio: 0.2,
	},
	{
		Category:   "food",
		IdealRatio: 0.2,
	},
	{
		Category:   "fixed",
		IdealRatio: 0.45,
	},
}

// TestCalculateScoreRule
func TestCalculateScoreRule(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name     string
		mockRepo func(m *mocks.MockRatioRuleRepository)
		req      *aggregate.TransactionReportDto
		want     float64
		wantErr  bool
		err      error
	}{
		{
			name: "リポジトリに接続時にエラーが発生した場合",
			mockRepo: func(m *mocks.MockRatioRuleRepository) {
				m.EXPECT().SelectAll(ctx).Return(nil, errors.New("mock err")).Times(1)
			},
			req: &aggregate.TransactionReportDto{
				BaseAmounts:   300000,
				TotalAmount:   200000,
				FoodExpenses:  20000,
				WasteExpenses: 10000,
				OtherExpenses: 10000,
				FixedCosts:    150000,
				VariableCosts: 10000,
				Savings:       100000,
			},
			want:    0,
			wantErr: true,
			err:     errors.New("mock err"),
		},
		{
			name: "スコアのバリデーションでエラーが発生した場合",
			mockRepo: func(m *mocks.MockRatioRuleRepository) {
				m.EXPECT().SelectAll(ctx).Return([]*models.RatioRule{
					{IdealRatio: 2.0},
				}, nil).Times(1)
			},
			req: &aggregate.TransactionReportDto{
				BaseAmounts:   300000,
				TotalAmount:   200000,
				FoodExpenses:  20000,
				WasteExpenses: 10000,
				OtherExpenses: 10000,
				FixedCosts:    150000,
				VariableCosts: 10000,
				Savings:       100000,
			},
			want:    0,
			wantErr: true,
			err:     errors.New("理想的な割合の合計が1.0ではありません"),
		},
		{
			name: "計算したスコアが0未満になった場合",
			mockRepo: func(m *mocks.MockRatioRuleRepository) {
				m.EXPECT().SelectAll(ctx).Return(normalRatioModels, nil).Times(1)
			},
			req: &aggregate.TransactionReportDto{
				BaseAmounts:   1000,
				TotalAmount:   1000000,
				FoodExpenses:  400000,
				WasteExpenses: 200000,
				OtherExpenses: 200000,
				FixedCosts:    200000,
				VariableCosts: 100000,
				Savings:       0,
			},
			want:    0,
			wantErr: true,
			err:     errors.New("スコアが0未満になりました"),
		},
		{
			name: "スコアが正常通り計算できた場合",
			mockRepo: func(m *mocks.MockRatioRuleRepository) {
				m.EXPECT().SelectAll(ctx).Return(normalRatioModels, nil).Times(1)
			},
			req: &aggregate.TransactionReportDto{
				BaseAmounts:   300000,
				TotalAmount:   200000,
				FoodExpenses:  20000,
				WasteExpenses: 10000,
				OtherExpenses: 10000,
				FixedCosts:    150000,
				VariableCosts: 10000,
				Savings:       100000,
			},
			want:    78.33333333333333,
			wantErr: false,
			err:     nil,
		},
	}

	for _, tt := range tests {
		t.Run("calculate score.", func(t *testing.T) {
			// mockの生成
			mock := gomock.NewController(t)
			defer mock.Finish()
			rrr := mocks.NewMockRatioRuleRepository(mock)
			//　mockの設定
			tt.mockRepo(rrr)
			// test対象の生成
			rule := NewCalculateScoreRule(rrr)

			// 計算処理
			got, err := rule.CalculateScore(ctx, tt.req)
			// errorの検証
			if (err != nil) != tt.wantErr {
				assert.Equal(t, tt.err.Error(), err.Error())
				// 通常の場合
			} else {
				assert.Equal(t, tt.want, got)
			}

		})
	}

}

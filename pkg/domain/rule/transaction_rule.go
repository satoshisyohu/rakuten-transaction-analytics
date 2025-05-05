//go:generate mockgen -source=transaction_rule.go -destination=../../../internal/mocks/domain/rule/transaction_rule.go -package=mocks

package rule

import "github.com/satoshisyohu/rakuten-transaction-analytics/pkg/domain/models"

type ITransactionRule interface {
	// CalculateAmount 使用した合計金額を計算する
	CalculateAmount() int64
	// CalculateFoodExpenses 食費を計算する(外食費用も含む)
	CalculateFoodExpenses() int64
	// CalculateWasteExpenses コンビニ費用を計算する
	CalculateWasteExpenses() int64
	// CalculateOtherExpenses 食費、コンビニ以外の費用を計算する
	CalculateOtherExpenses() int64
	// CalculateFixedCosts 固定費を計算する
	CalculateFixedCosts() int64
	// CalculateVariableCosts 変動費を計算する
	CalculateVariableCosts() int64
	// GetTransactions トランザクションを取得する
	GetTransactions() []*models.Transaction
	// Categorize カテゴライズする
	Categorize()
}

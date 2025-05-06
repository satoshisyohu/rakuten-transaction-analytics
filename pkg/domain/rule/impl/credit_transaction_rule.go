package impl

import (
	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/code"
	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/domain/models"
	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/domain/rule"
)

// CreditTransactionRule クレジットの取引明細を扱うルール
type CreditTransactionRule struct {
	creditTransaction []*models.Transaction
	cr                rule.ICategorizeRule
}

// NewCreditTransactionRule クレジットトランザクションルールを作成する
func NewCreditTransactionRule(transactions []*models.Transaction, cr rule.ICategorizeRule) rule.ITransactionRule {
	return &CreditTransactionRule{
		creditTransaction: transactions,
		cr:                cr,
	}
}

// CalculateAmount 使用した合計金額を計算する
func (ctr *CreditTransactionRule) CalculateAmount() int64 {
	var amount int64
	for _, v := range ctr.creditTransaction {
		amount += v.Amount
	}
	return amount
}

// CalculateFoodExpenses 食費を計算する(外食費用も含む)
func (ctr *CreditTransactionRule) CalculateFoodExpenses() int64 {
	var amount int64
	for _, v := range ctr.creditTransaction {
		if ctr.cr.IsFoodCosts(v.Detail) {
			amount += v.Amount
			// カテゴリーする
			v.Category = code.Food.String()
		}
	}
	return amount
}

// CalculateWasteExpenses 無駄遣いを計算する
func (ctr *CreditTransactionRule) CalculateWasteExpenses() int64 {
	var amount int64
	for _, v := range ctr.creditTransaction {

		if ctr.cr.IsWasteCosts(v.Detail) {
			//fmt.Println("コンビニ：", v.Detail, v.Amount)
			amount += v.Amount
			// カテゴリーする
			v.Category = code.Waste.String()
		}
	}
	return amount
}

// CalculateOtherExpenses 食費、コンビニ以外の費用を計算する
func (ctr *CreditTransactionRule) CalculateOtherExpenses() int64 {
	var amount int64
	for _, v := range ctr.creditTransaction {

		if ctr.cr.IsOtherCosts(v.Detail) {
			//fmt.Println("その他：", v.Detail, v.Amount)
			amount += v.Amount
			v.Category = code.Other.String()
		}
	}
	return amount
}

// CalculateFixedCosts 固定費を計算する
func (ctr *CreditTransactionRule) CalculateFixedCosts() int64 {
	var amount int64
	for _, v := range ctr.creditTransaction {

		if ctr.cr.IsContainFixedCosts(v.Detail) {
			amount += v.Amount
			v.Category = code.Fixed.String()
		}
	}
	return amount
}

// CalculateVariableCosts 変動費を計算する
func (ctr *CreditTransactionRule) CalculateVariableCosts() int64 {
	var amount int64
	for _, v := range ctr.creditTransaction {

		if ctr.cr.IsVariableCosts(v.Detail) {
			amount += v.Amount
			v.Category = code.Variable.String()
		}
	}
	return amount
}

// GetTransactions すべての明細を返却する
func (ctr *CreditTransactionRule) GetTransactions() []*models.Transaction {
	return ctr.creditTransaction
}

// Categorize カテゴライズする
func (ctr *CreditTransactionRule) Categorize() {
	for _, v := range ctr.creditTransaction {
		switch {
		case ctr.cr.IsFoodCosts(v.Detail):
			v.Category = code.Food.String()
		case ctr.cr.IsWasteCosts(v.Detail):
			v.Category = code.Waste.String()
		case ctr.cr.IsContainFixedCosts(v.Detail):
			v.Category = code.Fixed.String()
		case ctr.cr.IsVariableCosts(v.Detail):
			v.Category = code.Variable.String()
		case ctr.cr.IsOtherCosts(v.Detail):
			v.Category = code.Other.String()
		}
	}
}

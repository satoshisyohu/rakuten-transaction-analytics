package impl

import (
	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/code"
	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/domain/models"
	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/domain/rule"
)

// DebitTransactionRule デビットの取引明細を扱うルール
type DebitTransactionRule struct {
	debitTransaction []*models.Transaction
	cr               rule.ICategorizeRule
}

// NewDebitTransactionRule クレジットトランザクションルールを作成する
func NewDebitTransactionRule(transactions []*models.Transaction, cr rule.ICategorizeRule) rule.ITransactionRule {
	return &DebitTransactionRule{
		debitTransaction: transactions,
		cr:               cr,
	}
}

// CalculateAmount 使用した合計金額を計算する
func (dtr *DebitTransactionRule) CalculateAmount() int64 {
	var amount int64
	for _, v := range dtr.debitTransaction {
		amount += v.Amount
		v.Category = code.Food.String()
	}
	return amount
}

// CalculateFoodExpenses 食費を計算する(外食費用も含む)
func (dtr *DebitTransactionRule) CalculateFoodExpenses() int64 {
	var amount int64
	for _, v := range dtr.debitTransaction {
		if dtr.cr.IsFoodCosts(v.Detail) {
			//fmt.Println("食費：", v.Detail, v.Amount)
			amount += v.Amount
			// カテゴリーする
			v.Category = code.Food.String()
		}
	}
	return amount
}

// CalculateWasteExpenses 無駄遣いを計算する
func (dtr *DebitTransactionRule) CalculateWasteExpenses() int64 {
	var amount int64
	for _, v := range dtr.debitTransaction {

		if dtr.cr.IsWasteCosts(v.Detail) {
			//fmt.Println("コンビニ：", v.Detail, v.Amount)
			amount += v.Amount
			// カテゴリーする
			v.Category = code.Waste.String()
		}
	}
	return amount
}

// CalculateOtherExpenses 食費、コンビニ以外の費用を計算する
func (dtr *DebitTransactionRule) CalculateOtherExpenses() int64 {
	var amount int64
	for _, v := range dtr.debitTransaction {

		if dtr.cr.IsOtherCosts(v.Detail) {
			//fmt.Println("その他：", v.Detail, v.Amount)
			amount += v.Amount
			// カテゴリーする
			v.Category = code.Other.String()
		}
	}
	return amount
}

// CalculateFixedCosts 固定費を計算する
func (dtr *DebitTransactionRule) CalculateFixedCosts() int64 {
	var amount int64
	for _, v := range dtr.debitTransaction {

		if dtr.cr.IsContainFixedCosts(v.Detail) {
			amount += v.Amount
			v.Category = code.Fixed.String()
		}
	}
	return amount
}

// CalculateVariableCosts 変動費を計算する
func (dtr *DebitTransactionRule) CalculateVariableCosts() int64 {
	var amount int64
	for _, v := range dtr.debitTransaction {

		if dtr.cr.IsVariableCosts(v.Detail) {
			amount += v.Amount
			v.Category = code.Variable.String()
		}
	}
	return amount
}

// GetTransactions すべての明細を返却する
func (dtr *DebitTransactionRule) GetTransactions() []*models.Transaction {
	return dtr.debitTransaction

}

// Categorize カテゴライズする
func (dtr *DebitTransactionRule) Categorize() {
	for _, v := range dtr.debitTransaction {
		switch {
		case dtr.cr.IsFoodCosts(v.Detail):
			v.Category = code.Food.String()
		case dtr.cr.IsWasteCosts(v.Detail):
			v.Category = code.Waste.String()
		case dtr.cr.IsContainFixedCosts(v.Detail):
			v.Category = code.Fixed.String()
		case dtr.cr.IsVariableCosts(v.Detail):
			v.Category = code.Variable.String()
		case dtr.cr.IsOtherCosts(v.Detail):
			v.Category = code.Other.String()
		}
	}
}

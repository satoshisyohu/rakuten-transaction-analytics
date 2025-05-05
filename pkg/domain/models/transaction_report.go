package models

import (
	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/domain/aggregate"
)

type TransactionReport struct {
	// ID
	Id string
	// YearMonth
	YearMonth string
	// BaseAmount 基準金額
	BaseAmounts int64
	// 総使用金額
	TotalAmount int64
	// 食費
	FoodExpenses int64
	// 無駄遣い費用
	WasteExpenses int64
	// その他
	OtherExpenses int64
	// 固定費
	FixedCosts int64
	// 変動費
	VariableCosts int64
	//　スコア
	Score float64
	// 貯金
	Savings int64
}

func NewTransactionReport(yearMonth, id string, score float64, trd *aggregate.TransactionReportDto) *TransactionReport {
	return &TransactionReport{
		Id:            id,
		YearMonth:     yearMonth,
		BaseAmounts:   trd.BaseAmounts,
		TotalAmount:   trd.TotalAmount,
		FoodExpenses:  trd.FoodExpenses,
		WasteExpenses: trd.WasteExpenses,
		OtherExpenses: trd.OtherExpenses,
		FixedCosts:    trd.FixedCosts,
		VariableCosts: trd.VariableCosts,
		Savings:       trd.Savings,
		Score:         score,
	}
}

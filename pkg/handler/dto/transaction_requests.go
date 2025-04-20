package dto

type TransactionRequest struct {
	YearMonth string `json:"yearMonth"`
	// ベースとなる金額
	BaseAmounts int64 `json:"baseAmounts"`
	// 食費
	FoodExpenses int64 `json:"foodExpenses"`
	// コンビニ費用
	WasteExpenses int64 `json:"wasteExpenses"`
	// その他
	OtherExpenses int64 `json:"otherExpenses"`
	// 固定費
	FixedCosts int64 `json:"fixedCosts"`
	// 変動費
	VariableCosts int64 `json:"variableCosts"`
}

func (tr *TransactionRequest) CalculateTotalAmount() int64 {
	return tr.FoodExpenses + tr.WasteExpenses + tr.OtherExpenses + tr.FixedCosts + tr.VariableCosts
}

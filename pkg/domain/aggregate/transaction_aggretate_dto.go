package aggregate

// TransactionReportDto スコア算出時に使用するdto
type TransactionReportDto struct {
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
	// 貯金
	Savings int64
}

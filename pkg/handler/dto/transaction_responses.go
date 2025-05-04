package dto

import "cloud.google.com/go/civil"

type TransactionResponse struct {
	// 合計使用金額
	TotalAmount int64 `json:"totalAmount"`
	// 食費
	FoodExpenses int64 `json:"foodExpenses"`
	// コンビニ費用
	ConvenienceStoreExpenses int64 `json:"convenienceStoreExpenses"`
	// その他
	OtherExpenses int64 `json:"otherExpenses"`
	// 固定費
	FixedCosts int64 `json:"fixedCosts"`
	// 変動費
	VariableCosts int64 `json:"variableCosts"`
	// スコア
	Score float64 `json:"score"`
	// 貯金
	Savings int64 `json:"savings"`
}

type RetrieveTransactionResponses struct {
	Transactions []*RetrieveTransaction `json:"transactions"`
}

type RetrieveTransaction struct {
	TransactionDate civil.Date `json:"transactionDate"`
	Detail          string     `json:"detail"`
	Amount          int64      `json:"amount"`
	Category        string     `json:"category"`
	TransactionType string     `json:"transactionType"`
}

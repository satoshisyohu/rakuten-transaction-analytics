package code

// TransactionType TransactionTypeのコード値

type TransactionType string

const (
	// DebitTransaction DebitTransactionのコード値
	DebitTransaction TransactionType = "0"
	// CreditTransaction CreditTransactionのコード値
	CreditTransaction TransactionType = "1"
)

func (tt TransactionType) String() string {
	return string(tt)
}

// ToLabel ラベルに変換する
func (tt TransactionType) ToLabel() string {
	switch tt {
	case DebitTransaction:
		return "Debit"
	case CreditTransaction:
		return "Credit"
	default:
		return "Unknown"
	}
}

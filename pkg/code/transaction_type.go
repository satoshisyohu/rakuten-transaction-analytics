package code

type TransactionType string

const (
	DebitTransaction  TransactionType = "0"
	CreditTransaction TransactionType = "1"
)

func (tt TransactionType) String() string {
	return string(tt)
}

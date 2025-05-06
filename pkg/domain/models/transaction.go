package models

import (
	"time"

	"cloud.google.com/go/civil"
	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/code"
	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/helper"
)

// Transaction 明細情報
type Transaction struct {
	TransactionDate civil.Date
	Detail          string
	Amount          int64
	Category        string
	TransactionType string
	ID              string
}

// NewTransaction Transactionのファクトリ関数
func NewTransaction(date time.Time, storeName, transactionType string, amount int64) *Transaction {
	// idはscore生成時に設定するのでこの時点では空文字
	return &Transaction{
		TransactionDate: civil.Date{Year: date.Year(), Month: date.Month(), Day: date.Day()},
		Detail:          storeName,
		Amount:          amount,
		TransactionType: transactionType,
	}
}

// NewTransactions Transactionのスライスを作成するファクトリ関数、トランザクションタイプに応じてヘダーから値を取り出
func NewTransactions(records [][]string, transactionType code.TransactionType) []*Transaction {
	// recordsはtransactionの形式に変換して返却する
	var transaction []*Transaction
	for _, r := range records {
		timeDate := newTransactionDate(r, transactionType)
		amount := newAmount(r, transactionType)
		t := NewTransaction(timeDate, r[1], transactionType.ToLabel(), amount)
		transaction = append(transaction, t)
	}
	return transaction
}

func newTransactionDate(record []string, transactionType code.TransactionType) time.Time {
	switch transactionType {
	case code.DebitTransaction:
		return helper.StringToDate("20060102", record[0])
	case code.CreditTransaction:
		return helper.StringToDate("2006/01/02", record[0])
	default:
		// 不明なことは基本的にないがその場合変換せずに今日日付を返却する
		return time.Now()
	}
}
func newAmount(record []string, transactionType code.TransactionType) int64 {
	switch transactionType {
	case code.DebitTransaction:
		return helper.StringToInt64(record[2])
	case code.CreditTransaction:
		return helper.StringToInt64(record[4])
	default:
		// 不明なことは基本的にないがその場合0を返却する
		return 0
	}
}

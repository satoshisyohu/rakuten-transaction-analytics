package impl

import (
	"context"
	"errors"
	"mime/multipart"

	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/code"
	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/domain/models"
	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/domain/rule"
	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/domain/rule/impl"
	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/handler/dto"
)

// IRetrieveTransactionUsecase 取引明細を取得するユースケースの構造体
type RetrieveTransactionUsecase struct {
	fileReaderRule rule.IFileReaderRule
}

// NewRetrieveTransactionUsecase RetrieveTransactionUsecaseのファクトリ関数
func NewRetrieveTransactionUsecase(
	fileReaderRule rule.IFileReaderRule,
) *RetrieveTransactionUsecase {
	return &RetrieveTransactionUsecase{
		fileReaderRule: fileReaderRule,
	}
}

// Run ユースケースを実行する
func (r *RetrieveTransactionUsecase) Run(_ context.Context, files []*multipart.FileHeader) (*dto.RetrieveTransactionResponses, error) {

	var allTransaction []*models.Transaction
	// ファイルを読み込む
	// リクエストのファイル数分処理を行う
	for _, file := range files {
		records, err := r.fileReaderRule.ReadAll(file)
		if err != nil {
			return nil, err
		}
		// ヘダーの数に応じてruleを取得する
		transactionRule := r.handleTransactionRule(records)
		if transactionRule == nil {
			return nil, errors.New("transactionRuleが定義されていません。")
		}

		// カテゴライズする
		transactionRule.Categorize()

		// 明細をリストに詰める
		allTransaction = append(allTransaction, transactionRule.GetTransactions()...)
	}

	// レスポンスを返却する
	resTransactions := make([]*dto.RetrieveTransaction, len(allTransaction))
	for i, t := range allTransaction {
		resTransactions[i] = &dto.RetrieveTransaction{
			TransactionDate: t.TransactionDate,
			Detail:          t.Detail,
			Amount:          t.Amount,
			Category:        t.Category,
			TransactionType: t.TransactionType,
		}
	}

	return &dto.RetrieveTransactionResponses{
		Transactions: resTransactions,
	}, nil
}

// handleTransactionRule ヘダーの数に応じてトランザクションを返す（ヘダーの数が変わった際に影響を受けるので要注意）
func (r *RetrieveTransactionUsecase) handleTransactionRule(records [][]string) rule.ITransactionRule {
	// ヘダーの数に応じて対象のルールを判定する
	switch len(records[0]) {
	case CreditTransaction:
		return impl.NewCreditTransactionRule(models.NewTransactions(records[1:], code.CreditTransaction), &impl.CreditCategorizeRule{})
	case DebitTransaction:
		return impl.NewDebitTransactionRule(models.NewTransactions(records[1:], code.DebitTransaction), &impl.DebitCategorizeRule{})
	}
	// 基本的に想定外あるがルールが見つからない場合はnilを返す
	return nil
}

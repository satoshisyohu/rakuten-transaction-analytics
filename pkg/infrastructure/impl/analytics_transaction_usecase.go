package impl

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/code"
	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/domain/models"
	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/domain/repository"
	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/domain/rule"
	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/domain/rule/impl"
	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/handler/dto"
	"mime/multipart"
)

const (
	CreditTransaction = 10
	DebitTransaction  = 11
)

// AnalyticsTransactionUsecase トランザクションの分析を行うユースケース
type AnalyticsTransactionUsecase struct {
	fileReaderRule        rule.IFileReaderRule
	calculateScoreRule    rule.ICalculateScoreRule
	transactionRepository repository.TransactionRepository
	transactionReport     repository.TransactionReportRepository
}

// NewAnalyticsTransactionUsecase 新しいAnalyticsTransactionUsecaseを作成する
func NewAnalyticsTransactionUsecase(
	fileReaderRule rule.IFileReaderRule,
	calculateScoreRule rule.ICalculateScoreRule,
	transactionRepository repository.TransactionRepository,
	monthlySummaryRepository repository.TransactionReportRepository,
) *AnalyticsTransactionUsecase {
	return &AnalyticsTransactionUsecase{
		fileReaderRule:        fileReaderRule,
		calculateScoreRule:    calculateScoreRule,
		transactionRepository: transactionRepository,
		transactionReport:     monthlySummaryRepository,
	}
}

// Run ユースケースを実行する
func (a *AnalyticsTransactionUsecase) Run(ctx context.Context, req dto.TransactionRequest, files []*multipart.FileHeader) (*dto.TransactionResponse, error) {

	// リクエストの値をベースとして設定する。
	// 設定されていない場合は0となる
	u := uuid.New().String()
	transactionReport := &models.TransactionReport{
		Id:            u,
		YearMonth:     req.YearMonth,
		BaseAmounts:   req.BaseAmounts,
		TotalAmount:   req.CalculateTotalAmount(),
		FoodExpenses:  req.FoodExpenses,
		WasteExpenses: req.WasteExpenses,
		OtherExpenses: req.OtherExpenses,
		FixedCosts:    req.FixedCosts,
	}

	// 冪等性を担保するためにリクエストの月のTransactionReportsを取得する
	// レコードが存在する場合、基本的に処理は行わないが、リクエストのパラメタに応じて処理を行う
	tr, err := a.transactionReport.SelectByYearMonth(ctx, req.YearMonth)
	if err != nil {
		return nil, err
	}
	if len(tr) > 0 {
		// todo レコードが存在していてもパラメータに応じて削除して再度insertするようにする
		return nil, errors.New("すでに処理済みの月です。")
	}

	// すべての明細を溜め込むリスト
	var allTransaction []*models.Transaction

	// リクエストのファイル数分処理を行う
	for _, file := range files {
		var records [][]string
		records, err = a.fileReaderRule.ReadAll(file)
		if err != nil {
			return nil, err
		}
		// ヘダーの数に応じてruleを取得する
		transactionRule := a.handleTransactionRule(records, u)
		if transactionRule == nil {
			return nil, errors.New("transactionRuleが定義されていません。")
		}
		transactionReport.TotalAmount += transactionRule.CalculateAmount()
		transactionReport.FoodExpenses += transactionRule.CalculateFoodExpenses()
		transactionReport.WasteExpenses += transactionRule.CalculateWasteExpenses()
		transactionReport.FixedCosts += transactionRule.CalculateFixedCosts()
		transactionReport.VariableCosts += transactionRule.CalculateVariableCosts()
		transactionReport.OtherExpenses += transactionRule.CalculateOtherExpenses()

		// カテゴライズする
		transactionRule.Categorize()
		// 明細をリストに詰める
		allTransaction = append(allTransaction, transactionRule.GetTransactions()...)
	}
	// 貯金を計算する
	transactionReport.Savings = req.BaseAmounts - transactionReport.TotalAmount
	// スコアを計算する
	transactionReport.Score, err = a.calculateScoreRule.CalculateScore(transactionReport)
	if err != nil {
		return nil, err
	}

	if err = a.transactionRepository.SaveAll(ctx, allTransaction); err != nil {
		return nil, err
	}
	if err = a.transactionReport.SaveAll(ctx, []*models.TransactionReport{transactionReport}); err != nil {
		return nil, err
	}

	return &dto.TransactionResponse{
		// 合計使用金額
		TotalAmount: transactionReport.TotalAmount,
		// 食費
		FoodExpenses: transactionReport.FoodExpenses,
		// コンビニ費用
		ConvenienceStoreExpenses: transactionReport.WasteExpenses,
		// 固定費
		FixedCosts: transactionReport.FixedCosts,
		// 変動費
		VariableCosts: transactionReport.VariableCosts,
		// その他
		OtherExpenses: transactionReport.OtherExpenses,
		// 貯金
		Savings: transactionReport.Savings,
		// スコア
		Score: transactionReport.Score,
	}, nil

}

// handleTransactionRule ヘダーの数に応じてトランザクションを返す（ヘダーの数が変わった際に影響を受けるので要注意）
func (a *AnalyticsTransactionUsecase) handleTransactionRule(records [][]string, uuid string) rule.ITransactionRule {
	// ヘダーの数に応じて対象のルールを判定する
	switch len(records[0]) {
	case CreditTransaction:
		return impl.NewCreditTransactionRule(models.NewTransactions(records[1:], code.CreditTransaction, uuid), &impl.CreditCategorizeRule{})
	case DebitTransaction:
		return impl.NewDebitTransactionRule(models.NewTransactions(records[1:], code.DebitTransaction, uuid), &impl.DebitCategorizeRule{})
	}
	// ルールが見つからない場合はnilを返す
	return nil
}

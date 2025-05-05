package impl

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/code"
	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/domain/aggregate"
	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/domain/models"
	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/domain/repository"
	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/domain/rule"
	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/handler/dto"
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

func (a *AnalyticsTransactionUsecase) Run(ctx context.Context, req dto.TransactionRequest) (res *dto.TransactionResponse, err error) {
	// 冪等性を担保するためにリクエストの月のTransactionReportsを取得する
	//レコードが存在する場合、基本的に処理は行わないが、リクエストのパラメタに応じて処理を行う
	tr, err := a.transactionReport.SelectByYearMonth(ctx, req.YearMonth)
	if err != nil {
		return nil, err
	}
	if len(tr) > 0 {
		// todo レコードが存在していてもパラメータに応じて削除して再度insertするようにする
		return nil, errors.New("すでに処理済みの月です。")
	}

	var (
		// レコード投入時に付与するuuid
		u = uuid.New().String()
		// score
		score float64
	)
	// 分類ごとに当月の使用金額を計算し、DBに投入するためのレコードを作成する
	trd, modelTransaction := a.calculateAndAggregateTransaction(req, u)

	// スコアを計算する
	score, err = a.calculateScoreRule.CalculateScore(ctx, trd)
	if err != nil {
		return nil, err
	}

	// dtoとscoreからDBに保存するmodelを生成する
	transactionReport := &models.TransactionReport{
		Id:            u,
		YearMonth:     req.YearMonth,
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

	// すべての明細をDBに保存する
	if err = a.transactionRepository.SaveAll(ctx, modelTransaction); err != nil {
		return nil, err
	}
	// 明細のレポートをDBに保存する
	if err = a.transactionReport.SaveAll(ctx, []*models.TransactionReport{transactionReport}); err != nil {
		return nil, err
	}

	// responseの生成
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

// calculateAndAggregateTransaction ヘダーの数に応じてトランザクションを返す（ヘダーの数が変わった際に影響を受けるので要注意）
func (a *AnalyticsTransactionUsecase) calculateAndAggregateTransaction(req dto.TransactionRequest, u string) (*aggregate.TransactionReportDto, []*models.Transaction) {
	var (
		allTransaction = make([]*models.Transaction, len(req.Transactions))
		totalAmount    int64
		foodExpenses   = req.FoodExpenses
		wasteExpenses  = req.WasteExpenses
		fixedCosts     = req.FixedCosts
		variableCosts  = req.VariableCosts
		otherExpenses  = req.OtherExpenses
	)

	for i, t := range req.Transactions {
		// 分類に応じて金額を積み上げていく
		switch t.Category {
		case code.Other.String():
			otherExpenses += t.Amount
		case code.Food.String():
			foodExpenses += t.Amount
		case code.Waste.String():
			wasteExpenses += t.Amount
		case code.Fixed.String():
			fixedCosts += t.Amount
		case code.Variable.String():
			variableCosts += t.Amount
		default:
			// 何もしない
		}
		// dbに投入するためmodel形式に変換する
		transaction := &models.Transaction{
			TransactionDate: t.TransactionDate,
			Detail:          t.Detail,
			Amount:          t.Amount,
			Category:        t.Category,
			TransactionType: t.TransactionType,
			Id:              u,
		}
		allTransaction[i] = transaction
	}
	// リクエストにカテゴリごとの合計を計算した値を加算する
	totalAmount = a.getTotalAmount(foodExpenses, wasteExpenses, otherExpenses, fixedCosts, variableCosts)

	// score計算のためにdtoに詰める
	return &aggregate.TransactionReportDto{
		BaseAmounts:   req.BaseAmounts,
		TotalAmount:   totalAmount,
		FoodExpenses:  foodExpenses,
		WasteExpenses: wasteExpenses,
		OtherExpenses: otherExpenses,
		FixedCosts:    fixedCosts,
		VariableCosts: variableCosts,
		Savings:       req.BaseAmounts - totalAmount,
	}, allTransaction

}
func (a *AnalyticsTransactionUsecase) getTotalAmount(amounts ...int64) int64 {
	var total int64
	for _, v := range amounts {
		total += v
	}
	return total
}

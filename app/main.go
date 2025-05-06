package main

import (
	"cloud.google.com/go/bigquery"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/adapter/db"
	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/adapter/repository"
	rimpl "github.com/satoshisyohu/rakuten-transaction-analytics/pkg/domain/rule/impl"
	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/handler"
	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/helper"
	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/infrastructure/impl"
)

// main エントリーポイント
func main() {

	// todo pathは環境変数から読み込む
	err := helper.Load()
	if err != nil {
		fmt.Println(err)
		return
	}

	bqclient, err := db.NewBigqueryClient()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(bqclient *bigquery.Client) {
		dErr := bqclient.Close()
		if dErr != nil {
			slog.Error("Failed to close BigQuery client.", "err", dErr)
		}
	}(bqclient)

	// repository
	transactionRepository := repository.NewTransactionRepository(bqclient)
	monthlySummaryRepository := repository.NewMonthlySummaryRepository(bqclient)
	ratioRuleRepository := repository.NewRatioRuleRepository(bqclient)

	// rule
	fileReaderRule := rimpl.NewFileReaderRule()
	calculateScoreRule := rimpl.NewCalculateScoreRule(ratioRuleRepository)

	// usecase
	transactionUsecase := impl.NewAnalyticsTransactionUsecase(fileReaderRule, calculateScoreRule, transactionRepository, monthlySummaryRepository)
	retrieveTransactionUsecase := impl.NewRetrieveTransactionUsecase(fileReaderRule)

	// handler
	transactionHandler := handler.NewTransactionHandler(transactionUsecase, retrieveTransactionUsecase)

	http.HandleFunc("/rakuten/transaction/retrieve", transactionHandler.Retrieve)
	http.HandleFunc("/rakuten/transaction/execute", transactionHandler.Execute)

	fmt.Println("Server is running on http://localhost:8080")
	if err = http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println(err)
		return
	}

}

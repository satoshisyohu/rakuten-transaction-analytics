package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/adapter/repository"
	rimpl "github.com/satoshisyohu/rakuten-transaction-analytics/pkg/domain/rule/impl"
	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/handler"
	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/helper"
	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/infrastructure/impl"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
	"net/http"
	"os"
)

func main() {

	// todo pathは環境変数から読み込む
	err := helper.Load()
	if err != nil {
		fmt.Println(err)
		return
	}

	bqclient, err := repository.NewBigqueryClient()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer bqclient.Close()

	// repository
	transactionRepository := repository.NewTransactionRepository(bqclient)
	monthlySummaryRepository := repository.NewMonthlySummaryRepository(bqclient)

	// rule
	fileReaderRule := rimpl.NewFileReaderRule()
	calculateScoreRule := rimpl.NewCalculateScoreRule()
	// usecase
	transactionUsecase := impl.NewAnalyticsTransactionUsecase(fileReaderRule, calculateScoreRule, transactionRepository, monthlySummaryRepository)
	// handler
	transactionHandler := handler.NewTransactionHandler(transactionUsecase)

	http.HandleFunc("/rakuten/transaction/analysis", transactionHandler.Analysis)
	fmt.Println("Server is running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println(err)
		return
	}
	//usecase := impl.AnalyticsTransactionUsecase{}

	//for _, f := range []string{"debit.csv", "credit.csv"} {
	//	file, err := os.Open(f)
	//	// // エラーチェック
	//	if err != nil {
	//		panic(err)
	//	}
	//	defer file.Close()
	//	// CSVリーダーを作成
	//	//reader := csv.NewReader(transform.NewReader(file, japanese.ShiftJIS.NewDecoder()))
	//
	//	reader, err := NewReader(file)
	//	if err != nil {
	//		panic(err)
	//	}
	//	// CSVを読み込む
	//	records, err := reader.ReadAll()
	//	// エラーチェック
	//	if err != nil {
	//		panic(err)
	//	}
	//	usecase.Run(records)
	//}

}

func NewReader(file *os.File) (*csv.Reader, error) {

	// BOM判定のためファイルの先頭から3バイトを読み込む
	bom := make([]byte, 3)
	_, err := file.Read(bom)
	if err != nil {
		return nil, err
	}

	// bomの判定ロジックを実施する
	expectedBOM := []byte{0xEF, 0xBB, 0xBF}
	// UTF-8のBOMか判定する
	if bytes.Equal(bom, expectedBOM) {
		// BOMが存在する場合はすでにファイルポインタがBOM以降のためそのまま返す
		// これはクレジットカードの明細を読み込んだ場合でありUTF-8のためそのまま読み込み可能
		return csv.NewReader(file), nil
	} else {
		// BOMがない場合、ファイルポインタを先頭に戻す必要がある
		// これはデビットカードの明細を読み込んだ場合でありShift-JISのためデコードして読み込む必要がある
		_, err = file.Seek(0, 0)
		if err != nil {
			return nil, err
		}
		// BOMがない場合はShift-JISでデコードする
		return csv.NewReader(transform.NewReader(file, japanese.ShiftJIS.NewDecoder())), nil
	}
}

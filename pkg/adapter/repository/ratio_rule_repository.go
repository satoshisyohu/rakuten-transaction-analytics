package repository

import (
	"context"
	"errors"
	"fmt"
	"log"

	"cloud.google.com/go/bigquery"
	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/domain/models"
	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/domain/repository"
	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/helper"
	"google.golang.org/api/iterator"
)

const (
	ratioRuleTable = "RatioRules"
)

// transactionRepository 明細を保存するリポジトリ
type ratioRuleRepository struct {
	client *bigquery.Client
}

// NewRatioRuleRepository コンストラクタ
func NewRatioRuleRepository(client *bigquery.Client) repository.RatioRuleRepository {
	return &ratioRuleRepository{
		client: client,
	}
}

// SelectAll 明細を保存する
func (rr *ratioRuleRepository) SelectAll(ctx context.Context) ([]*models.RatioRule, error) {
	sql := fmt.Sprintf(
		"SELECT Id,IdealRatio,Category,StartDate,EndDate FROM `%s.%s`",
		helper.GetDatasetId(), ratioRuleTable,
	)

	it, err := rr.client.Query(sql).Read(ctx)
	if err != nil {
		fmt.Println("Error reading query results:", err)
		return nil, err
	}
	var res []*models.RatioRule
	for {
		var r models.RatioRule
		err := it.Next(&r)
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate result: %v", err)
		}
		res = append(res, &r)
	}
	return res, nil

}

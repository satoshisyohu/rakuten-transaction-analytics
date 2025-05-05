package repository

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"cloud.google.com/go/bigquery"
	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/domain/models"
	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/domain/repository"
	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/helper"
	"google.golang.org/api/iterator"
)

const (
	ratioRuleTable = "RatioRules"
	selectAllSql   = "SELECT Id,IdealRatio,Category,StartDate,EndDate FROM `%s.%s`"
)

// ratioRuleRepository 明細を保存するリポジトリ
type ratioRuleRepository struct {
	client *bigquery.Client
}

// NewRatioRuleRepository コンストラクタ
func NewRatioRuleRepository(client *bigquery.Client) repository.RatioRuleRepository {
	return &ratioRuleRepository{
		client: client,
	}
}

// SelectAll 比率のルールを取得する
func (rr *ratioRuleRepository) SelectAll(ctx context.Context) ([]*models.RatioRule, error) {
	sql := fmt.Sprintf(selectAllSql, helper.GetDatasetId(), ratioRuleTable)

	it, err := rr.client.Query(sql).Read(ctx)
	if err != nil {
		slog.Error("Error reading query results.", "err", err)
		return nil, err
	}
	var res []*models.RatioRule
	for {
		var r models.RatioRule
		err = it.Next(&r)
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			slog.Error("Failed to iterate result.", "err", err)
		}
		res = append(res, &r)
	}
	return res, nil

}

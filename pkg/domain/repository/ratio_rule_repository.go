package repository

import (
	"context"

	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/domain/models"
)

type RatioRuleRepository interface {
	SelectAll(ctx context.Context) ([]*models.RatioRule, error)
}

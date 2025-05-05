//go:generate mockgen -source=ratio_rule_repository.go -destination=../../../internal/mocks/domain/repository/ratio_rule_repository.go -package=mocks
package repository

import (
	"context"

	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/domain/models"
)

type RatioRuleRepository interface {
	SelectAll(ctx context.Context) ([]*models.RatioRule, error)
}

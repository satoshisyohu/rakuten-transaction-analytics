package models

import "github.com/satoshisyohu/rakuten-transaction-analytics/pkg/code"

type Score struct {
	IdealRatio  float64
	ActualRatio float64
	Category    code.Category
}

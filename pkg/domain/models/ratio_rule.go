package models

import (
	"cloud.google.com/go/civil"
)

// RatioRule 比率のマスタデータを管理する構造体
type RatioRule struct {
	ID         string
	IdealRatio float64
	Category   string
	StartDate  civil.Date
	EndDate    civil.Date
}

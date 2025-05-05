package models

import (
	"cloud.google.com/go/civil"
)

type RatioRule struct {
	Id         string
	IdealRatio float64
	Category   string
	StartDate  civil.Date
	EndDate    civil.Date
}

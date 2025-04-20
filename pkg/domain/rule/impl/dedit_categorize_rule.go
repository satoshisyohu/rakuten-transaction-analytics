package impl

import (
	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/helper"
)

type DebitCategorizeRule struct {
}

// IsFoodCosts 食費
func (d *DebitCategorizeRule) IsFoodCosts(val string) bool {
	return d.isContainSuperMarketCosts(val) || d.isContainEatingOutCosts(val)
}

// IsContainSuperMarketCosts スーパーの判定
func (d *DebitCategorizeRule) isContainSuperMarketCosts(val string) bool {
	return helper.IsContain(val, []string{
		"マックスバリュ",
		"にしてつストア",
		"ＭＡＲＫ　ＩＳ福岡ももち",
		"ﾕﾒﾏ-ﾄｸﾏﾓﾄｻﾆ-",
		"やまや",   // パスタを買うところ
		"ハローデイ", // マークイズももち
		"天神地下街", // 怪しいが、カルディとして食費にカウントする
	})
}

// IsContainEatingOutCosts 外食費の判定
func (d *DebitCategorizeRule) isContainEatingOutCosts(val string) bool {
	return helper.IsContain(val, []string{"玄海丸"})
}

// IsWasteCosts 無駄遣いの判定
func (d *DebitCategorizeRule) IsWasteCosts(val string) bool {
	return helper.IsContain(val, []string{"セブン－イレブン", "ローソン"})
}

// IsContainFixedCosts 固定費の判定
func (d *DebitCategorizeRule) IsContainFixedCosts(_ string) bool {
	// 現状は該当なし
	return false
}

// IsVariableCosts 変動費の判定
func (d *DebitCategorizeRule) IsVariableCosts(val string) bool {
	return helper.IsContain(val, []string{"きゅうでんガス", "九州電力"})
}

// IsOtherCosts その他の費用の判定
func (d *DebitCategorizeRule) IsOtherCosts(val string) bool {
	return !d.IsFoodCosts(val) && !d.IsWasteCosts(val) &&
		!d.IsContainFixedCosts(val) && !d.IsVariableCosts(val)
}

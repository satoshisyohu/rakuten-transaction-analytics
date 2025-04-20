package impl

import "github.com/satoshisyohu/rakuten-transaction-analytics/pkg/helper"

type CreditCategorizeRule struct {
}

// IsFoodCosts 食費
func (c *CreditCategorizeRule) IsFoodCosts(val string) bool {
	return c.isContainSuperMarketCosts(val) || c.isContainEatingOutCosts(val)
}

// isContainSuperMarketCosts スーパーの判定
func (c *CreditCategorizeRule) isContainSuperMarketCosts(val string) bool {
	return helper.IsContain(val, []string{
		"マックスバリュ", //todo 未確認
		"西鉄ストア",
		"ＭＡＲＫ　ＩＳ福岡ももち", //todo 未確認
		"サニー",
		"やまや",   // todo 未確認 パスタを買うところ
		"ハローデイ", // todo 未確認マークイズももち
		"天神地下街", //  //todo 未確認 怪しいが、カルディとして食費にカウントする
	})
}

// isContainEatingOutCosts 外食費の判定
func (c *CreditCategorizeRule) isContainEatingOutCosts(val string) bool {
	return helper.IsContain(val, []string{
		"玄海丸", //todo 未確認
	})
}

// IsWasteCosts 無駄遣いの判定
func (c *CreditCategorizeRule) IsWasteCosts(val string) bool {
	return helper.IsContain(val, []string{
		"セブン－イレブン", // todo 未確認
		"ローソン",     // todo 未確認
	})
}

// IsContainFixedCosts 固定費の判定
func (c *CreditCategorizeRule) IsContainFixedCosts(val string) bool {
	return helper.IsContain(val, []string{"ソフトバンク"})
}

// IsVariableCosts 変動費の判定
func (c *CreditCategorizeRule) IsVariableCosts(val string) bool {
	return helper.IsContain(val, []string{"きゅうでんガス", "九州電力"})
}

// IsOtherCosts その他の費用の判定
func (c *CreditCategorizeRule) IsOtherCosts(val string) bool {
	return !c.IsFoodCosts(val) && !c.IsWasteCosts(val) && !c.IsContainFixedCosts(val) && !c.IsVariableCosts(val)
}

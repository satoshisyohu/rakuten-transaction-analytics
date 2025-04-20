package rule

type ICategorizeRule interface {
	// IsFoodCosts　食費
	IsFoodCosts(val string) bool
	// IsWasteCosts 無駄遣い
	IsWasteCosts(val string) bool
	// IsContainFixedCosts　変動費
	IsContainFixedCosts(val string) bool
	// IsVariableCosts　固定費
	IsVariableCosts(val string) bool
	// IsOtherCosts　その他
	IsOtherCosts(val string) bool
}

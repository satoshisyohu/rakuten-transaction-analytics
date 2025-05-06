//go:generate mockgen -source=categorize_rule.go -destination=../../../internal/mocks/domain/rule/categorize_rule.go -package=mocks

package rule

// ICategorizeRule interface
type ICategorizeRule interface {
	// IsFoodCosts　食費
	IsFoodCosts(string) bool
	// IsWasteCosts 無駄遣い
	IsWasteCosts(string) bool
	// IsContainFixedCosts　変動費
	IsContainFixedCosts(string) bool
	// IsVariableCosts　固定費
	IsVariableCosts(string) bool
	// IsOtherCosts　その他
	IsOtherCosts(string) bool
}

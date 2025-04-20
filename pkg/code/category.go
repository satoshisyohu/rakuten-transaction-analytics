package code

type Category string

const (
	Food     Category = "food"
	Other    Category = "other"
	Waste    Category = "waste"
	Fixed    Category = "fixed"
	Variable Category = "variable"
	Savings  Category = "variable"
)

func (c Category) String() string {
	return string(c)
}

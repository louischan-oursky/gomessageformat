//go:generate stringer -type Complexity $GOFILE
package messageformat

type Complexity int

const (
	NoComplexity Complexity = iota
	SingleLiteral
	SingleRoad
	Complex
)

func (c Complexity) IsSingleLiteral() bool { return c == SingleLiteral }
func (c Complexity) IsSingleRoad() bool    { return c == SingleRoad }
func (c Complexity) IsComplex() bool       { return c == Complex }

package messageformat

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/empirefox/makeplural/plural"
	"golang.org/x/text/language"
)

var Funcs = template.FuncMap{
	"mfPluralNode":         ParsePlural,
	"mfPluralNodeWithData": ParsePluralWithData,
}

func ParsePlural(culture language.Tag, input string) (Node, error) {
	return ParsePluralWithData(culture, input, nil)
}

func ParsePluralWithData(culture language.Tag, input string, data interface{}) (Node, error) {
	_, lang, ok := plural.Info.Find(culture)
	if !ok {
		return nil, fmt.Errorf("plural culture not found: %s", culture)
	}

	p, err := NewWithCulture(lang)
	if err != nil {
		return nil, err
	}

	mf, err := p.Parse(input, data)
	if err != nil {
		return nil, err
	}

	return mf.Root(), nil
}

type Node interface {
	Type() string
	Culture() language.Tag
	Complexity() Complexity
	Data() interface{}
	Root() Node
	Parent() Node
	Children() []Node
	Format(mf *MessageFormat, buf *bytes.Buffer, data Data, pound string) error
	setParent(parent Node)
	addChild(child Node)
}

type Root struct {
	Container
	culture    language.Tag
	complexity Complexity
	data       interface{}
}

func (nd *Root) Type() string          { return "root" }
func (nd *Root) Culture() language.Tag { return nd.culture }
func (nd *Root) Data() interface{}     { return nd.data }
func (nd *Root) Root() Node            { return nd }
func (nd *Root) setParent(parent Node) { panic("can not call Root.setParent") }

func (nd *Root) Complexity() Complexity {
	if nd.complexity == NoComplexity {
		if nd.isSingleLiteral() {
			nd.complexity = SingleLiteral
		} else if alwaysSingleChild(nd) {
			nd.complexity = SingleRoad
		} else {
			nd.complexity = Complex
		}
	}
	return nd.complexity
}

func (nd *Root) isSingleLiteral() bool {
	switch len(nd.Children()) {
	case 0:
		return true
	case 1:
		if ltr, ok := nd.Children()[0].(*Literal); ok && isSingleLiteral(ltr) {
			return true
		}
	default:
	}
	return false
}

func isSingleLiteral(nd *Literal) bool { return len(nd.Content) < 2 }

func alwaysSingleChild(nd Node) bool {
	if _, ok := nd.(SelectNode); ok {
		for _, c := range nd.Children() {
			if !alwaysSingleChild(c) {
				return false
			}
		}
		return true
	}

	if ltr, ok := nd.(*Literal); ok {
		return isSingleLiteral(ltr)
	}

	if len(nd.Children()) == 0 {
		return true
	}

	if len(nd.Children()) != 1 {
		return false
	}

	return alwaysSingleChild(nd.Children()[0])
}

type Container struct {
	root     Node
	parent   Node
	children []Node
}

func newContainer(parent Node) *Container {
	child := new(Container)
	AddChild(parent, child)
	return child
}

func (nd *Container) Type() string           { return "container" }
func (nd *Container) Culture() language.Tag  { return nd.Root().Culture() }
func (nd *Container) Complexity() Complexity { return nd.Root().Complexity() }
func (nd *Container) Data() interface{}      { return nd.Root().Data() }
func (nd *Container) Root() Node             { return nd.root }
func (nd *Container) Parent() Node           { return nd.parent }
func (nd *Container) Children() []Node       { return nd.children }
func (nd *Container) setParent(parent Node) {
	nd.root = parent.Root()
	nd.parent = parent
}
func (nd *Container) addChild(child Node) {
	nd.children = append(nd.children, child)
}

func AddChild(parent, child Node) {
	child.setParent(parent)
	parent.addChild(child)
}

func (nd *Container) Format(mf *MessageFormat, output *bytes.Buffer, data Data, pound string) error {
	for _, child := range nd.children {
		err := child.Format(mf, output, data, pound)
		if err != nil {
			return err
		}
	}
	return nil
}

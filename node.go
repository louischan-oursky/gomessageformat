package messageformat

import (
	"bytes"
	"text/template"

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
	p, err := NewWithCulture(culture)
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
	Data() interface{}
	Root() Node
	Parent() Node
	Children() []Node
	Add(child Node)
	Format(mf *MessageFormat, buf *bytes.Buffer, data Data, pound string) error
	setRoot(root Node)
	setParent(parent Node)
}

type Root struct {
	Container
	culture language.Tag
	data    interface{}
}

func (nd *Root) Type() string          { return "root" }
func (nd *Root) Culture() language.Tag { return nd.culture }
func (nd *Root) Data() interface{}     { return nd.data }
func (nd *Root) Root() Node            { return nd }
func (nd *Root) setRoot(root Node)     { panic("can not call Root.setRoot") }
func (nd *Root) setParent(parent Node) { panic("can not call Root.setParent") }

type Container struct {
	root     Node
	parent   Node
	children []Node
}

func newContainer(parent Node) *Container {
	child := new(Container)
	parent.Add(child)
	return child
}

func (nd *Container) Type() string          { return "container" }
func (nd *Container) Culture() language.Tag { return nd.Root().Culture() }
func (nd *Container) Data() interface{}     { return nd.Root().Data() }
func (nd *Container) Root() Node            { return nd.root }
func (nd *Container) Parent() Node          { return nd.parent }
func (nd *Container) Children() []Node      { return nd.children }
func (nd *Container) setRoot(root Node)     { nd.root = root }
func (nd *Container) setParent(parent Node) { nd.parent = parent }

func (parent *Container) Add(child Node) {
	child.setRoot(parent.Root())
	child.setParent(parent)
	parent.children = append(parent.children, child)
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

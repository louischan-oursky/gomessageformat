package messageformat

import (
	"bytes"
	"text/template"
)

var Funcs = template.FuncMap{
	"mfPluralNode": ParsePlural,
}

func ParsePlural(cultrue, input string) (Node, error) {
	p, err := NewWithCulture(cultrue)
	if err != nil {
		return nil, err
	}

	mf, err := p.Parse(input)
	if err != nil {
		return nil, err
	}

	return mf.Root(), nil
}

type Node interface {
	Type() string
	Culture() string
	Root() Node
	Parent() Node
	Children() []Node
	Add(child Node)
	Format(mf *MessageFormat, buf *bytes.Buffer, data Data, pound string) error
	setRoot(root Node)
	setParent(parent Node)
}

type Container struct {
	culture  string
	root     Node
	parent   Node
	children []Node
}

func newContainer(parent Node) *Container {
	child := new(Container)
	parent.Add(child)
	return child
}

func (nd *Container) Type() string { return "container" }
func (nd *Container) Culture() string {
	if nd.culture != "" {
		return nd.culture
	}
	return nd.Root().Culture()
}
func (nd *Container) Root() Node {
	if nd.root == nil {
		return nd
	}
	return nd.root
}
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

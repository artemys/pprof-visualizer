package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"sort"
	"strings"
)

//go:embed tree.html
var TreeNodeTemplate string

type FunctionsTree struct {
	Name string    `json:"name"`
	Root *TreeNode `json:"root"`
}

func (t *FunctionsTree) sort() {
	if t.Root == nil {
		fmt.Println("warn: called sort() on an empty tree")
		return
	}

	t.Root.sort()
}

type TreeNode struct {
	Children []*TreeNode `json:"itemId"`
	Function Function    `json:"function"`
	Self     int64       `json:"self"`
	Value    int64       `json:"value"`
	Percent  float64     `json:"percent"`
	Visible  bool        `json:"visible"`
}

func NewFunctionsTree(treeName string) *FunctionsTree {
	return &FunctionsTree{
		Name: treeName,
		Root: &TreeNode{},
	}
}

func (n TreeNode) ID(lineNumber bool) string {
	if n.Function.Name == "" {
		return "Root"
	}
	return n.Function.String(lineNumber)
}

// AddFunction adds the given Function to the tree.
// AddFunction takes care of aggregating the values per functions calls or line of
// code depending on the aggregateByFunction parameter.
func (n *TreeNode) AddFunction(f Function, value, self int64, percent float64, aggregateByFunction bool) *TreeNode {
	for i, child := range n.Children {
		// if existing, we add the values to the current node
		if child.ID(!aggregateByFunction) == f.String(!aggregateByFunction) {
			child.Value += value
			child.Self += self
			child.Percent += percent
			n.Children[i] = child
			return child
		}
	}

	// doesn't exist, create it
	node := &TreeNode{
		Function: f,
		Value:    value,
		Self:     self,
		Percent:  percent,
	}

	n.Children = append(n.Children, node)
	return node
}

func (n *TreeNode) isLeaf() bool {
	return len(n.Children) == 0
}

func (n *TreeNode) filter(searchField string) bool {
	var visible bool

	if searchField == "" || n.Function.Name == "" {
		visible = true
	} else if strings.Contains(strings.ToLower(n.Function.Name), strings.ToLower(searchField)) {
		visible = true
	} else if strings.Contains(strings.ToLower(n.Function.File), strings.ToLower(searchField)) {
		visible = true
	}

	for _, child := range n.Children {
		if child.filter(searchField) {
			visible = true
		}
	}

	n.Visible = visible
	return n.Visible
}

func (n *TreeNode) sort() {
	sort.Slice(
		n.Children,
		func(i, j int) bool {
			return n.Children[i].Value > n.Children[j].Value
		},
	)
	for _, child := range n.Children {
		child.sort()
	}
}

func (t *FunctionsTree) toHtml() string {
	var b bytes.Buffer
	b.WriteString("<div>")
	b.WriteString(fmt.Sprintf("<p>%s</p>", t.Name))
	b.WriteString("<ul>")
	b.WriteString(t.Root.toHtml())
	b.WriteString("</ul>")
	b.WriteString("</div>")
	return b.String()
}
func (t *TreeNode) toHtml() string {
	var b bytes.Buffer
	b.WriteString("<li>")
	var buf bytes.Buffer
	tmpl, _ := template.New("tree").Parse(TreeNodeTemplate)
	tmpl.Execute(&buf, t)
	b.WriteString(buf.String())
	if t.Children != nil && len(t.Children) > 0 {
		b.WriteString("<ul>")
		for _, child := range t.Children {
			b.WriteString(child.toHtml())
		}
		b.WriteString("</ul>")
	}
	b.WriteString("</li>")
	return b.String()
}

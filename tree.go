package main

import (
	"fmt"
	"sort"
	"strings"
)

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

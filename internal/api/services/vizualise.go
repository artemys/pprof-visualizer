package services

import (
	"fmt"
	"github.com/artemys/pprof-visualizer/internal/pkg/pprof"
	"path"
	"time"
)

type Result struct {
	TotalAllocBytes int64 `json:"total_alloc_bytes"`
}

func Visualize(profile pprof.Profile) *FunctionsTree {
	p, _ := NewProfile(&profile, "")
	ftree := p.BuildTree("tree", true, "")
	return ftree
}

func (p *Profile) texts(node *TreeNode) (value string, self string, tooltip string, lineText string) {
	if p.Type == "cpu" {
		value = time.Duration(node.Value).String()
		self = time.Duration(node.Self).String()
		tooltip = fmt.Sprintf("%s of %s\nself: %s", value, time.Duration(p.TotalSampling).String(), self)
	} else {

		value = humanize.IBytes(uint64(node.Value))
		self = humanize.IBytes(uint64(node.Self))
		tooltip = fmt.Sprintf("%s of %s\nself: %s", value, humanize.IBytes(p.TotalSampling), self)
	}
	lineText = fmt.Sprintf("%s %s:%d - %s - self: %s", node.function.Name, path.Base(node.function.File), node.function.LineNumber, value, self)
	if p.aggregateByFunction {
		lineText = fmt.Sprintf("%s %s - %s - self: %s", node.function.Name, path.Base(node.function.File), value, self)
	}
	return value, self, tooltip, lineText
}

package services

import (
	"fmt"
	"github.com/artemys/pprof-visualizer/internal/pkg/pprof"
	"github.com/dustin/go-humanize"
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
	lineText = fmt.Sprintf("%s %s:%d - %s - self: %s", node.Function.Name, path.Base(node.Function.File), node.Function.LineNumber, value, self)
	if p.aggregateByFunction {
		lineText = fmt.Sprintf("%s %s - %s - self: %s", node.Function.Name, path.Base(node.Function.File), value, self)
	}
	return value, self, tooltip, lineText
}

func (p *Profile) SetResume(name string) {
	switch p.Type {
	case ModeCpu:
		p.resume = fmt.Sprintf("%s - total sampling duration: %s - total capture duration %s", name, time.Duration(p.TotalSampling).String(), p.CaptureDuration.String())
	case ModeHeapAlloc:
		p.resume = fmt.Sprintf("%s - total allocated memory: %s", p., humanize.IBytes(p.TotalSampling))
	case ModeHeapInuse:
		p.resume = fmt.Sprintf("%s - total in-use memory: %s", name, humanize.IBytes(p.TotalSampling))
	default:
		p.resume = "fail"
	}
}

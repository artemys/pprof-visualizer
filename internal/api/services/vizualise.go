package services

import (
	"compress/gzip"
	"fmt"
	"github.com/artemys/pprof-visualizer/internal/pkg/pprof"
	"github.com/gogo/protobuf/proto"
	"io/ioutil"
	"os"
)

type Result struct {
	TotalAllocBytes int64 `json:"total_alloc_bytes"`
}

func Visualize(profile pprof.Profile) *FunctionsTree {
	p, _ := NewProfile(&profile, "")
	ftree := p.BuildTree("tree", true, "")
	return ftree
}

func readProtoFile(filename string) (*pprof.Profile, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("readProtoFile: os.Open: %v", err)
	}

	g, err := gzip.NewReader(f)
	if err != nil {
		return nil, fmt.Errorf("readProtoFile: gzip.NewReader: %v", err)
	}

	data, err := ioutil.ReadAll(g)
	if err != nil {
		return nil, fmt.Errorf("readProtoFile: ioutil.ReadAll: %v", err)
	}

	var profile pprof.Profile
	if err := proto.Unmarshal(data, &profile); err != nil {
		return nil, fmt.Errorf("readProtoFile: proto.Unmarshal: %v", err)
	}

	return &profile, nil
}

package services

import (
	"fmt"
	"github.com/artemys/pprof-visualizer/internal/pkg/logger"
	"github.com/artemys/pprof-visualizer/internal/pkg/pprof"
	"go.uber.org/zap"
	"os"
	"time"
)

const (
	// use this when you don't really know the mode
	// to use to read the profile.
	ModeDefault         string = ""
	ModeCpu, CpuKeyword string = "cpu", "cpu"
	ModeHeapAlloc       string = "heap-alloc"
	ModeHeapInuse       string = "heap-inuse"
	SpaceKeyword        string = "space"
)

type Profile struct {
	Samples
	TotalSampling   uint64
	CaptureDuration time.Duration

	// "cpu" or "heap"
	Type string
	Mode string

	functionsMapByLocation ManyFunctionsMap
	locationsMap           LocationsMap
	stringsMap             StringsMap
	aggregateByFunction    bool //nolint:unused
	resume                 string
	Name                   string
}

func NewProfile(p *pprof.Profile, mode string) (*Profile, error) {
	// start by building some maps because everything
	// is indexed in various maps.
	// ----------------------

	// strings map
	stringsMap := buildStringsTable(p)

	// functions map
	functionsMap := buildFunctionsMap(p, stringsMap)

	// locations map
	locationsMap, functionsMapByLocation := buildLocationsMap(p, functionsMap)

	// let's now build the profile
	// ----------------------

	typ := ReadProfileType(p)

	if typ != CpuKeyword && typ != SpaceKeyword {
		return nil, fmt.Errorf("unsupported type: %s", typ)
	}

	profile := readProfile(p, stringsMap, functionsMapByLocation, locationsMap, mode)

	switch typ {
	case CpuKeyword:
		profile.Type = CpuKeyword
	case SpaceKeyword:
		profile.Type = "heap"
	}
	profile.SetResume("test")
	return profile, nil
}

func ReadProfileType(p *pprof.Profile) string {
	return p.StringTable[uint64(p.GetPeriodType().Type)]
}

func readProfile(p *pprof.Profile, stringsMap StringsMap, functionsMapByLocation ManyFunctionsMap,
	locationsMap LocationsMap, mode string) *Profile {
	var samples Samples
	var idx int

	switch {
	case mode == ModeDefault:
		fallthrough
	case ReadProfileType(p) == CpuKeyword && mode == ModeCpu:
		idx = 1
	case ReadProfileType(p) == SpaceKeyword && mode == ModeHeapAlloc:
		idx = 1
	case ReadProfileType(p) == SpaceKeyword && mode == ModeHeapInuse:
		idx = 3
	default:
		logger.Log.Info("err: incompatible mode and profile type.",
			zap.String("readProfileType", ReadProfileType(p)),
			zap.String("mode", mode))
		os.Exit(-1)
	}

	for _, pprofSample := range p.Sample {
		var sample Sample
		// cpu [1] cpu usage
		// space [1] heap allocated
		// space [3] heap in use
		value := pprofSample.GetValue()[idx]

		for i := len(pprofSample.LocationId) - 1; i >= 0; i-- {
			l := pprofSample.LocationId[i]
			sample.Functions = append(sample.Functions, functionsMapByLocation[l]...)
			sample.Value = value
		}

		// compute the Self time for the leaf
		leaf := sample.Functions[len(sample.Functions)-1]
		leaf.Self += value
		sample.Functions[len(sample.Functions)-1] = leaf

		samples = append(samples, sample)
	}

	// compute the total sampling time
	var totalSum uint64
	for _, s := range samples {
		totalSum += uint64(s.Value)
	}

	// compute the percentage for every sample
	for i, s := range samples {
		s.PercentTotal = float64(s.Value) / (float64(totalSum)) * 100.0
		samples[i] = s
	}

	return &Profile{
		Samples:         samples,
		TotalSampling:   totalSum,
		CaptureDuration: time.Duration(p.GetDurationNanos()),

		functionsMapByLocation: functionsMapByLocation,
		locationsMap:           locationsMap,
		stringsMap:             stringsMap,
		Name:                   "test",
	}
}

func (p *Profile) BuildTree(treeName string, aggregateByFunction bool, searchField string) *FunctionsTree {
	tree := NewFunctionsTree(treeName)

	for _, s := range p.Samples {
		node := tree.Root
		for _, f := range s.Functions {
			if s.Value == 0 {
				continue
			}
			node = node.AddFunction(f, s.Value, f.Self, s.PercentTotal, aggregateByFunction)
		}
	}

	if tree.Root != nil {
		tree.Root.filter(searchField)
	}

	tree.sort()

	return tree
}

func buildLocationsMap(profile *pprof.Profile,
	functionsMap FunctionsMap) (LocationsMap, ManyFunctionsMap) {
	rv := make(LocationsMap)
	lrv := make(ManyFunctionsMap)

	for _, location := range profile.Location {
		if location.Line[0] == nil {
			continue
		}

		loc := Location{}

		for idx := len(location.Line) - 1; idx >= 0; idx-- {
			line := location.Line[idx]
			inlined := idx != len(location.Line)-1

			f := functionsMap[line.GetFunctionId()]
			f.LineNumber = uint64(line.GetLine())
			if inlined {
				f.Name = fmt.Sprintf("(inlined) %s", f.Name)
			}
			loc.Functions = append(loc.Functions, f)

			// set the line number in functions map if not inlined
			if !inlined {
				f := functionsMap[line.GetFunctionId()]
				f.LineNumber = uint64(line.GetLine())
				functionsMap[line.GetFunctionId()] = f
			}

			fs := lrv[location.GetId()]
			lrv[location.GetId()] = append(fs, f)

			rv[location.GetId()] = loc
		}
	}

	return rv, lrv
}

func buildFunctionsMap(profile *pprof.Profile, stringsMap StringsMap) FunctionsMap {
	rv := make(FunctionsMap)
	for _, f := range profile.Function {
		rv[f.GetId()] = Function{
			Name: stringsMap[uint64(f.GetName())],
			File: stringsMap[uint64(f.GetFilename())],
		}
	}
	return rv
}

func buildStringsTable(profile *pprof.Profile) StringsMap {
	rv := make(StringsMap)
	for i, v := range profile.GetStringTable() {
		rv[uint64(i)] = v
	}
	return rv
}

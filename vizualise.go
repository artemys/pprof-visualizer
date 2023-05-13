package main

import (
	"compress/gzip"
	"encoding/json"
	"example.com/pprof-visualizer/pprof"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"io/ioutil"
	"net/http"
	"os"
)

type Result struct {
	TotalAllocBytes int64 `json:"total_alloc_bytes"`
}

func Visualize(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer gzipReader.Close()

	// Lecture du contenu décompressé du fichier "pprof.pb"
	data, err := ioutil.ReadAll(gzipReader)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var profile pprof.Profile
	if err := proto.Unmarshal(data, &profile); err != nil {
		fmt.Println("error reading file")
		return
	}
	p, _ := NewProfile(&profile, "")
	ftree := p.BuildTree("tree", true, "")
	fmt.Println(ftree.root)

	response := ftree
	jsonData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
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

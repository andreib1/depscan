package main

import (
	"encoding/json"
	"github.com/crystal-construct/depscan/analyze"
	"github.com/crystal-construct/depscan/gomod"
	"github.com/crystal-construct/depscan/graphviz"
	"io/ioutil"
	"os"
	"path"
	"sort"
)

func main() {
	prj := os.Args[1]
	f1, _ := os.Open(path.Join(prj, "go.sum"))
	f2, _ := os.Open(path.Join(prj, "go.mod"))
	defer f1.Close()
	defer f2.Close()
	modMap := &gomod.ModMap{}
	gomod.Parse(f2, f1, modMap)
	b, err := json.MarshalIndent(modMap, "", "  ")
	if err != nil {
		panic(err)
	}
	ioutil.WriteFile("deps.json", b, 0777)
	graphviz.Create("gv.dot", modMap)
	vers := analyze.Scan(modMap)
	f, err := os.OpenFile("deps.txt", os.O_RDWR+os.O_CREATE+os.O_TRUNC, 0777)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	for i, j := range vers {
		f.Write([]byte(i))
		sort.Strings(j)
		for _, k := range j {
			f.Write([]byte("\t" + k + "\n"))
		}
	}
}

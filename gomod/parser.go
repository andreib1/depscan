package gomod

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
	"regexp"
	"strings"
)

type Mod struct {
	Name         string
	Version      string
	Hash         string
	ModHash      string
	Comment      string `json:",omitempty"`
	Deps         map[string]*Mod `json:",omitempty"`
}

func (m * Mod) Copy() *Mod {
	return &Mod{
		Name:m.Name,
		Version:m.Version,
		Hash:m.Hash,
		ModHash:m.ModHash,
		Comment:m.Comment,
		Deps: make(map[string]*Mod),
	}
}

type ModMap struct {
	Root *Mod
	Mods map[string] *Mod `json:""`
}

var sumrx *regexp.Regexp
var modrx *regexp.Regexp
var modProbe *regexp.Regexp
var gopath string
func init() {
	sumrx, _ = regexp.Compile("^([^\\ \\t]*)\\s(v[0-9]*\\.[0-9]*\\.[0-9]*[^\\s]*)\\s?(h1:)?([^\\ ]*)?")
	modrx, _ = regexp.Compile(`^\s*([^\s^\t]*)\s+(v[0-9]+\.[0-9]+\.[0-9]+[^\s]*)\s*(h1:)?([^\s]*)`)
	gopath = path.Join(strings.Replace(os.Getenv("GOPATH"),"\\", "/",-1),"pkg","mod", "cache", "download")
	fmt.Println(gopath)
}

func Parse(mod io.ReadCloser, sum io.ReadCloser, modMap *ModMap) {
	parseGoSum(sum, modMap)
	modMap.Root = &Mod{}
	parseMod(mod, modMap.Root, modMap)
}

func parseMod(m io.ReadCloser, mod *Mod, modMap *ModMap) {
	if mod.Deps == nil {
		mod.Deps = make(map[string]*Mod)
	}
	s := bufio.NewScanner(m)
	for s.Scan() != false {
		if o := modrx.Match(s.Bytes()); o {
			mt := modrx.FindSubmatch(s.Bytes())
			en := string(mt[1])
			n := en + " " + string(mt[2])
			if strings.HasPrefix(en,"\"") && strings.HasSuffix(en,"\"") {
				en = strings.Trim(en,"\"")
			}
			n = en + " " + string(mt[2])
			en = convert(string(en))
			fmt.Println(n)
			mx := modMap.Mods[n].Copy()
			f,err := os.Open(path.Join(gopath,string(en),"@v",string(mt[2])+".mod"))
			if err != nil {
				panic(err)
			}
			defer f.Close()
			parseMod(f,mx, modMap)
			mod.Deps[n] = mx
		}
	}
}

func convert(s string) string {
	a := byte('A')
	z := byte('Z')
	st := []byte(s)
	st2 := []byte{}
	for i := 0; i < len(st); i++ {
		if (st[i] <= z && st[i] >= a) {
			st2 = append(st2, byte('!'), strings.ToLower(string([]byte{st[i]}))[0])
		} else {
			st2 = append(st2, st[i])
		}
	}
	return string(st2)
}

func parseGoSum(sum io.ReadCloser, modMap *ModMap) {
	if modMap.Mods == nil {
		modMap.Mods = make(map[string]*Mod)
	}
	s := bufio.NewScanner(sum)
	for s.Scan() {
		x := sumrx.FindSubmatch(s.Bytes())
		n := string(x[1])
		v := string(x[2])
		h := string(x[4])
		var mh bool
		if strings.HasSuffix(v,"/go.mod") {
			v = strings.TrimSuffix(v, "/go.mod")
			mh = true
		}
		var m *Mod
		if t, ok := modMap.Mods[n + " " + v]; ok {
			m = t
		} else {
			m = &Mod{
				Deps: make(map[string]*Mod),
			}
			modMap.Mods[n + " " + v] = m
		}
		if mh {
			v = strings.TrimSuffix(v, "/go.mod")
			m.Name = n
			m.Version = v
			m.ModHash = h
		} else {
			m.Name = n
			m.Version = v
			m.Hash = h
		}
	}
}
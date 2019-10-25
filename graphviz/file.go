package graphviz

import (
	"fmt"
	"github.com/crystal-construct/depscan/gomod"
	"os"
	"strings"

	"strconv"
)

func Create(filename string, modMap *gomod.ModMap) {
	s := make(map[string]string)
	n := 1
	f, err := os.OpenFile(filename, os.O_RDWR + os.O_CREATE + os.O_TRUNC, 0777)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.Write([]byte("digraph D {\n"))
	recurse(modMap.Root, s, f, &n)
	f.Write([]byte("}\n"))
}

func recurse(m *gomod.Mod, s map[string]string, f *os.File, n *int) {
	un := m.Name + " " + m.Version
	e, nw := see(un, s, n)
	if nw {

		for _, i := range m.Deps {
			recurse(i, s, f, n)
			f.Write([]byte(fmt.Sprintf("\t%s -> %s\n", e , s[i.Name + " " + i.Version])))
		}
		f.Write([]byte(fmt.Sprintf("\t%s [label=\"%s\"];\n", e, strings.Replace(un, " ", "\\n", -1))))
	}
}

func see(un string, s map[string]string, n *int) (string, bool) {
	var nw = false
	e := "a" + strconv.FormatInt(int64(*n),10)
	if q, ok := s[un]; ok {
		e = q
	} else {
		s[un] = e
		*n++
		nw = true
	}
	return e, nw
}
package analyze

import "github.com/crystal-construct/depscan/gomod"

func Scan(modMap *gomod.ModMap) map[string][]string {
	s := make(map[string][]string)
	for _, m := range modMap.Mods {
		if _, ok := s[m.Name]; !ok {
			s[m.Name] = make([]string,1)
		}
		s[m.Name] = append(s[m.Name], m.Version)
	}
	return s
}
package wordlists

import (
	"fmt"
	"strings"
)

//go:generate go run process_wordlist/main.go --name en --output eff_large_wordlist.go eff_large_wordlist.txt
//go:generate gofmt -w eff_large_wordlist.go

// Lists contains the list of valid word lists.
var Lists = make(map[string][]string)

// Get returns the word list name.
func Get(name string) ([]string, error) {
	if list, ok := Lists[name]; ok {
		return list, nil
	}

	return nil, fmt.Errorf("wordlist %q not found, valid wordlists: %v", name, strings.Join(Names(), ", "))
}

// Names returns the list of names of all available wordlists.
func Names() []string {
	names := make([]string, 0, len(Lists))
	for name := range Lists {
		names = append(names, name)
	}
	return names
}

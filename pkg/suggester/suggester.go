package suggester

import (
	"fmt"
	"strings"
	"sync"
)

// Suggester represents a sug suggestion finder
type Suggester struct {
	tlds sync.Map

	// options
	withList []string
}

// New initiates a new Suggester
func New(opts ...Option) *Suggester {
	sug := Suggester{
		tlds:     sync.Map{},
		withList: defaultTLDs[:],
	}

	sug.Options(opts...)
	sug.init()

	return &sug
}

// Options applies the given opts onto Suggester
func (sug *Suggester) Options(opts ...Option) {
	for _, opt := range opts {
		opt(sug)
	}
}

func (sug *Suggester) init() {
	for _, tld := range sug.withList {
		tld = strings.ToLower(tld)
		if skipTLD(tld) {
			continue
		}
		sug.tlds.Store(string(tld), true)
	}
}

// Suggest returns potential domains for the given input
// Output is a map where keys are potential tlds and values are the full potential domain name
func (sug *Suggester) Suggest(input ...string) (map[string]string, error) {
	suggestions := make(map[string]string)
	for _, in := range input {
		if len(in) < 4 {
			continue
		}

		tld := strings.ToLower(in)
		tld = in[len(in)-2:]
		secondaryLevel := in[:len(in)-len(tld)]
		for len(secondaryLevel) >= 2 {
			if sug.isExistingTLD(tld) {
				suggestions[tld] = fmt.Sprintf("%s.%s", secondaryLevel, tld)
			}

			tld = in[len(in)-(len(tld)+1):]
			secondaryLevel = in[:len(in)-len(tld)]
		}

	}

	return suggestions, nil
}

func (sug *Suggester) isExistingTLD(potentialTLD string) bool {
	_, ok := sug.tlds.Load(potentialTLD)
	return ok
}

// Option is a type for Suggester Options
type Option func(sug *Suggester)

// OptionWithList allows specifying a list for suggestions overriding the default TLDs list
func OptionWithList(list []string) Option {
	return func(sug *Suggester) {
		sug.withList = list
	}
}

func skipTLD(tld string) bool {
	// filtering xn--??? TLD
	return strings.HasPrefix(tld, "xn--")
}

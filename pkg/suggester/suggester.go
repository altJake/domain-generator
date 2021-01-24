package suggester

import (
	"fmt"
	"strings"
	"sync"

	"github.com/pkg/errors"
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
func (sug *Suggester) Suggest(input string) (map[string]string, error) {
	if len(input) < 4 {
		return nil, errors.New("Top Level Domain cannot contain less than 2 letters, structure must be ab.cd at the least")
	}

	input = strings.ToLower(input)
	collection := make(map[string]string)

	tld := input[len(input)-2:]
	secondaryLevel := input[:len(input)-len(tld)]
	for len(secondaryLevel) >= 2 {
		if sug.isExistingTLD(tld) {
			collection[tld] = fmt.Sprintf("%s.%s", secondaryLevel, tld)
		}

		tld = input[len(input)-(len(tld)+1):]
		secondaryLevel = input[:len(input)-len(tld)]
	}

	return collection, nil
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

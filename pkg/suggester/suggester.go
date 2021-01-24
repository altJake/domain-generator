package suggester

import (
	"fmt"
	"strings"
	"sync"
)

// Suggester represents a sug suggestion finder
type Suggester struct {
	tlds        map[string]bool
	wg          *sync.WaitGroup
	suggestions sync.Map

	// options
	withList []string
}

// New initiates a new Suggester
func New(opts ...Option) *Suggester {
	sug := Suggester{
		tlds:     map[string]bool{},
		wg:       &sync.WaitGroup{},
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
		sug.tlds[tld] = true
	}
}

// Suggest returns potential domains for the given input
// Output is a map where keys are potential tlds and values are the full potential domain name
func (sug *Suggester) Suggest(input ...string) (map[string]string, error) {
	sug.suggestions = sync.Map{}

	for _, in := range input {
		sug.wg.Add(1)
		go sug.processSingleInput(in)
	}

	sug.wg.Wait()

	result := make(map[string]string)
	sug.suggestions.Range(func(key, value interface{}) bool {
		tld := key.(string)
		domain := value.(string)
		result[tld] = domain

		return true
	})

	return result, nil
}

func (sug *Suggester) processSingleInput(input string) {
	defer sug.wg.Done()

	if len(input) < 4 {
		// skip when domain is invalid due to minimal legnth
		return
	}

	input = strings.ToLower(input)
	tld := input[len(input)-2:]
	secondaryLevel := input[:len(input)-len(tld)]

	for len(secondaryLevel) >= 2 {
		if sug.isExistingTLD(tld) {
			sug.suggestions.Store(tld, fmt.Sprintf("%s.%s", secondaryLevel, tld))
		}

		tld = input[len(input)-(len(tld)+1):]
		secondaryLevel = input[:len(input)-len(tld)]
	}
}

func (sug *Suggester) isExistingTLD(potentialTLD string) bool {
	_, ok := sug.tlds[potentialTLD]
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

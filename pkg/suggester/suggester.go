package suggester

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

// Suggester represents a sug suggestion finder
type Suggester struct {
	tlds sync.Map
}

// New initiates a new Suggester
func New() *Suggester {
	sug := Suggester{
		tlds: sync.Map{},
	}

	err := sug.load()
	if err != nil {
		panic(errors.Wrap(err, "Failed to load domains"))
	}

	return &sug
}

func (sug *Suggester) load() error {
	var err error
	filePath := getFilePath()

	dat, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	splits := bytes.Split(dat, []byte("\n"))
	for _, tld := range splits {
		tld = bytes.ToLower(tld)
		if invalidValidTLDFileLine(tld) {
			continue
		}
		sug.tlds.Store(string(tld), true)
	}

	return nil
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

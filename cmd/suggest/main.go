package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/altjake/domain-generator/pkg/suggester"

	"github.com/pkg/errors"
)

func main() {
	input, err := parse()
	if err != nil {
		fmt.Println(errors.Wrap(err, "Failed to parse command line arguments"))
		os.Exit(1)
	}

	suggester := suggester.New()
	domains, err := suggester.Suggest(input...)
	if err != nil {
		fmt.Println(errors.Wrap(err, "Execution Failed"))
		os.Exit(1)
	}

	output := getSafeOutput(domains)

	fmt.Printf("Potential domains for the given input: \n%s\n", output)
}

func parse() ([]string, error) {
	if len(os.Args) == 1 {
		return nil, errors.New("Missing command line argument as input")
	}
	return os.Args[1:], nil
}

func getSafeOutput(domains map[string]string) string {
	output, err := json.MarshalIndent(domains, "", "  ")
	if err != nil {
		return fmt.Sprint(domains)
	}

	return string(output)
}

package suggester

import (
	"bytes"
	"os"
)

const (
	defaultFilePath string = "resources/domains.txt"
)

func getFilePath() string {
	if envFilePath := os.Getenv("DOMGEN_FILE_PATH"); envFilePath != "" {
		return envFilePath
	}

	return defaultFilePath
}

func invalidValidTLDFileLine(tld []byte) bool {
	return bytes.Contains(tld, []byte("#")) || // commented line
		bytes.HasPrefix(tld, []byte("xn--")) // filtering xn--??? TLD
}

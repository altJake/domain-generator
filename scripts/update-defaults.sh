#!/bin/sh

OUTPUT_FILE=pkg/suggester/defaults.go

# validate necessary tools
command -v curl >/dev/null 2>&1 || { echo >&2 '"curl" is required to run this script.\n'; exit 1; }
command -v sed >/dev/null 2>&1 || { echo >&2 '"sed" is required to run this script.\n'; exit 1; }

echo '
╔════════════════════════════════════════════════╗
║        Downloading latest domains file         ║
╚════════════════════════════════════════════════╝'
DOMAINS=$(curl https://data.iana.org/TLD/tlds-alpha-by-domain.txt)

echo '
╔════════════════════════════════════════════════╗
║                Sanitizing file                 ║
╚════════════════════════════════════════════════╝'
DOMAINS=$(echo "$DOMAINS" | sed '/^#/d')

echo '
╔════════════════════════════════════════════════╗
║                 Generate file                  ║
╚════════════════════════════════════════════════╝'
cat <<EOT > $OUTPUT_FILE
//⚠️⚠️⚠️ THIS FILE IS GENERATED - Last updated at $(date -u +"%Y-%m-%d %H:%M:%S") ⚠️⚠️⚠️

package suggester

var defaultTLDs = [...]string{
	$(echo "$DOMAINS" | sed 's/.*/"&",/')
}

EOT

echo '
╔════════════════════════════════════════════════╗
║                  Format file                   ║
╚════════════════════════════════════════════════╝'
# silencing redundant output
gofmt -l -w $OUTPUT_FILE >/dev/null 

echo "\n$OUTPUT_FILE updated successfully!"

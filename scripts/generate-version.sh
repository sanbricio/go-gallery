#!/bin/bash
# scripts/generate-version.sh

VERSION=$(git describe --tags --always)
COMMIT=$(git rev-parse HEAD)
DATE=$(date -u +%Y-%m-%dT%H:%M:%SZ)

mkdir -p ../src/commons/configurator/version

cat <<EOF > ../src/commons/configurator/version/version.go
package version

var (
    AppVersion = "$VERSION"
    CommitHash = "$COMMIT"
    BuildDate  = "$DATE"
)
EOF

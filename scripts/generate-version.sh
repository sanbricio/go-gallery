#!/bin/bash
# scripts/generate-version.sh
MAIN_GO_FILE="../main.go"
BRANCH_NAME=$(git symbolic-ref --short HEAD)

echo "Branch: $BRANCH_NAME"
if [[ ! "$BRANCH_NAME" =~ ^release/v([0-9]+\.[0-9]+\.[0-9]+)$ ]]; then
    echo "Not a valid release branch. Exiting."
    exit 1
fi
VERSION="v${BASH_REMATCH[1]}"
echo "Release version: $VERSION"
COMMIT=$(git rev-parse HEAD)
DATE=$(date +'%d/%m/%Y %H:%M:%S')

# Create version.go file with current values
mkdir -p ../src/commons/configurator/version

cat <<EOF >../src/commons/configurator/version/version.go
package version

var (
    AppVersion = "$VERSION"
    CommitHash = "$COMMIT"
    BuildDate  = "$DATE"
)
EOF

echo "✔ Version generated: $VERSION"
echo "Updating the version in the file $MAIN_GO_FILE..."

# Update the version in main.go file
grep -q "^// @version " "$MAIN_GO_FILE"
sed -i "s|^// @version .*|// @version $VERSION|" "$MAIN_GO_FILE"
echo "✔ Version updated to $VERSION in main.go"

# Check if swag is installed
if ! command -v swag &> /dev/null; then
    echo "⚠ 'swag' not found. Installing..."
    go install github.com/swaggo/swag/cmd/swag@latest
    export PATH="$PATH:$(go env GOPATH)/bin"
    echo "✔ 'swag' installed successfully"
fi

# Generate Swagger documentation
echo "✔ Generating Swagger documentation..."
cd ..
swag init
echo "✔ Documentation generated successfully"

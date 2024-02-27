#!/usr/bin/env zsh

RELEASE="$1"
CURRENT_HASH=""

print_hash() {
    OS="$1"
    ARCH="$2"
    VERSION="$3"

    if [[ "$ARCH" == "x86_64" ]]; then
        EXT="${OS}_amd64.zip"
    else
        EXT="${OS}_arm64.zip"
    fi

    URL="https://github.com/databricks/cli/releases/download/v${VERSION}/databricks_cli_${VERSION}_${EXT}"

    CURRENT_HASH=$(nix store prefetch-file "$URL" --json | jq -r '.hash')

    echo "${VERSION}  ${OS}_${ARCH}: $CURRENT_HASH"
}

print_hash "linux" "x86_64" "$RELEASE"
print_hash "linux" "arm64" "$RELEASE"
print_hash "darwin" "x86_64" "$RELEASE"
print_hash "darwin" "arm64" "$RELEASE"
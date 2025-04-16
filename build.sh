#!/bin/bash

set -e

# Output directory
OUTPUT_DIR="./bin"
mkdir -p "$OUTPUT_DIR"

# Build configurations
BUILD_TARGETS=(
    "linux:amd64"
    "darwin:amd64"
    # "windows:amd64"
)

# Build loop
for TARGET in "${BUILD_TARGETS[@]}"; do
    OS=${TARGET%%:*}
    ARCH=${TARGET##*:}
    OUTPUT_NAME="fakerfactory-${OS}-${ARCH}"

    if [ "$OS" == "windows" ]; then
        OUTPUT_NAME+=".exe"
    fi

    echo "Building for $OS/$ARCH..."
    if [ "$OS" == "linux" ]; then
        export CC="x86_64-linux-musl-gcc"
        export CXX="x86_64-linux-musl-g++"
    fi
    if [ "$OS" == "darwin" ]; then
        export CC=""
        export CXX=""
    fi
    CGO_ENABLED=1 GOARCH=$ARCH GOOS=$OS go build -o "$OUTPUT_DIR/$OUTPUT_NAME" .
done

echo "Build completed. Binaries are in the $OUTPUT_DIR directory."
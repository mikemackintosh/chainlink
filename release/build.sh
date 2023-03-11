#!/bin/bash

source ./release/common.sh

# ARCHITECTURES SUPPORTED
ARCH=(
    amd64
    arm64
)

# OPERATING SYSTEMS SUPPORT
OS=(
    linux
    darwin
    windows
)

# STEP 2: Build the ldflags

LDFLAGS=(
  "-X '${PACKAGE}/version.Version=${VERSION}'"
  "-X '${PACKAGE}/version.CommitHash=${SHORT_COMMIT_HASH}'"
  "-X '${PACKAGE}/version.BuildTime=${BUILD_TIMESTAMP}'"
)

# STEP 3: Actual Go build process

go build -ldflags="${LDFLAGS[*]}"

if [[ ! -d release/bin ]]; then 
    mkdir release/bin
fi

for os in ${OS[@]}; do 
    for arch in ${ARCH[@]}; do
        echo -e "\nBuilding for $os-$arch"
        GOOS=$os GOARCH=$arch go build -ldflags="-s -w ${LDFLAGS[*]}" -o release/bin/ninetails-$os-$arch cmd/main.go
        echo $(shasum -a 256 release/bin/ninetails-$os-$arch)
    done
done

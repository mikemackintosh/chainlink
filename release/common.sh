export REPO="mikemackintosh/chainlink"
export PACKAGE="github.com/${REPO}/internal"
export VERSION="v$(git describe --tags --always --abbrev=0 --match='v[0-9]*.[0-9]*.[0-9]*' 2> /dev/null | sed 's/^.//')"
export COMMIT_HASH="$(git rev-parse HEAD)"
export SHORT_COMMIT_HASH="$(git rev-parse --short HEAD)"
export BUILD_TIMESTAMP=$(date '+%Y-%m-%dT%H:%M:%S')

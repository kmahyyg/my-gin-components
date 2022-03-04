#!/usr/bin/env zsh

CURVER=$(git describe --tags --always --dirty)
CURENV=$1
PROJECT_NAME=$(basename $(pwd))

echo "Current Project: ${PROJECT_NAME}"

GC_BUILDFLAGS="-ldflags \"-X cmd.versionNum=${CURVER} -s -w \" -trimpath"
GARBLE_BUILDFLAGS="-literals -seed=random"
SRC_FILES=$(find . -name "*.go" | grev -v "vendor")

export GOPROXY=https://goproxy.cn,direct
export CGO_ENABLED=0

PLATFORMS="linux/amd64"
ADDITIONAL_PLATFORMS="darwin/amd64 windows/386 linux/386 windows/amd64 linux/arm64 windows/arm64"
OUTPUT="${PROJECT_NAME}"

gofmt -s -w .

ALWAYS_BUILD() {
  for PLATFORM in $PLATFORMS; do
    GOOS=${PLATFORM%/*}
    GOARCH=${PLATFORM#*/}
    BIN_FILENAME="${OUTPUT}-${GOOS}-${GOARCH}"
    if [[ "${GOOS}" == "windows" ]]; then BIN_FILENAME="${BIN_FILENAME}.exe"; fi
    CMD="GOOS=${GOOS} GOARCH=${GOARCH} go build -o bin/${BIN_FILENAME} ${GC_BUILDFLAGS} ${SRC_FILES}"
    echo "${CMD}"
    eval $CMD
  done
}

CURRENT_MACHINE_BUILD() {
  go build -o bin/${OUTPUT}-$(uname -s)-$(uname -m) ${GC_BUILDFLAGS} ${SRC_FILES}
}

case $CURENV in
  "dev")
    echo "Building for development"
    GC_BUILDFLAGS="${GC_BUILDFLAGS} -gcflags \"-N -l\""
    ALWAYS_BUILD
    CURRENT_MACHINE_BUILD
    ;;

  "prod")
    echo "Building for production"
    PLATFORMS="$PLATFORMS $ADDITIONAL_PLATFORMS"
    ALWAYS_BUILD
    ;;

  *)
    echo "Unknown environment"
    exit 1
    ;;
esac

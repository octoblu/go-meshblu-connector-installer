#!/bin/bash

APP_NAME='meshblu-connector-assembler'
BUILD_DIR="$PWD/dist"
IMAGE_NAME="local/$APP_NAME"

build_on_docker() {
  docker build --tag "$IMAGE_NAME:built" .
}

build_on_local() {
  local goos="$1"
  local goarch="$2"

  env GOOS="$goos" GOARCH="$goarch" go build -a -tags netgo -installsuffix cgo -ldflags '-w' -o "${BUILD_DIR}/${APP_NAME}-${goos}-${goarch}" .
}

build_osx_on_local() {
  build_on_local 'darwin' 'amd64'
}

copy() {
  cp $BUILD_DIR/$APP_NAME* entrypoint/
}

init() {
  rm -rf "$BUILD_DIR/" \
    && mkdir -p "$BUILD_DIR/"
}

package() {
  docker build --tag "$IMAGE_NAME:latest" entrypoint
}

run() {
  docker run --rm \
    --volume "$BUILD_DIR:/export/" \
    "$IMAGE_NAME:built" \
      cp "$APP_NAME" '/export'
}

fatal() {
  local message="$1"
  echo "$message"
  exit 1
}

cross_compile_build(){
  for goos in darwin linux windows; do
    for goarch in 386 amd64 arm; do
      if [ "${goos}-${goarch}" == 'windows-arm' ]; then
        echo '* skipping windows-arm'
        continue
      fi
      echo "* building: ${goos}-${goarch}"
      build_on_local "$goos" "$goarch" > /dev/null
    done
  done
}

docker_build() {
  init            || fatal 'init failed'
  build_on_docker || fatal 'build_on_docker failed'
  run             || fatal 'run failed'
  copy            || fatal 'copy failed'
  package         || fatal 'package failed'
}

osx_build() {
  init               || fatal 'init failed'
  build_osx_on_local || fatal 'build_osx_on_local failed'
}

release_osx_build() {
  mkdir -p dist \
  && osx_build \
  && tar -czf "${APP_NAME}-osx.tar.gz" "${APP_NAME}" \
  && mv "${APP_NAME}-osx.tar.gz" dist/
  echo "* wrote dist/${APP_NAME}-osx.tar.gz"
}

main() {
  local mode="$1"
  if [ "$mode" == "docker" ]; then
    echo '# Docker Build'
    docker_build
    exit $?
  fi

  if [ "$mode" == 'osx' ]; then
    echo '# OSX Build'
    osx_build
    exit $?
  fi

  if [ "$mode" == 'release-osx' ]; then
    echo '# Release Build'
    release_osx_build
    exit $?
  fi

  if [ "$mode" == 'cross-compile' ]; then
    echo '# Cross Compile'
    cross_compile_build
    exit $?
  fi

  echo "Usage: ./build.sh docker/osx/release-osx/cross-compile"
  exit 1
}
main $@

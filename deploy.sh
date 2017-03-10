#!/bin/bash

APP_NAME=go-meshblu-connector-installer
LOCAL_IMAGE_NAME="local/$APP_NAME"
REMOTE_IMAGE_NAME="octoblu/$APP_NAME"

fatal() {
  local message=$1
  echo "$message"
  exit 1
}

main() {
  echo "#######################################################################"
  echo "# WARNING: You may now gump and codefresh will deploy automagically!! #"
  echo "#######################################################################"
  echo 
  read -p "Are you sure you want to manually deploy? " -n 1 -r
  echo    # (optional) move to a new line
  if [[ $REPLY =~ ^[Yy]$ ]]
  then
    local tag="$1"
    if [ -z "$tag" ]; then
      fatal "Missing tag as first argument"
    fi
    echo "building $LOCAL_IMAGE_NAME:$tag"
    ./build.sh "$tag" || fatal 'build failed'
    echo "tagging $LOCAL_IMAGE_NAME:$tag -> $REMOTE_IMAGE_NAME:$tag"
    docker tag "$LOCAL_IMAGE_NAME:$tag" "$REMOTE_IMAGE_NAME:$tag" || fatal 'failed to tag'
    echo "pushing $REMOTE_IMAGE_NAME:$tag"
    docker push "$REMOTE_IMAGE_NAME:$tag" || fatal 'failed to push'
  fi
}

main "$@"

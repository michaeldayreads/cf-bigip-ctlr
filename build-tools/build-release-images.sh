#!/bin/bash

# This script packages up artifacts produced by ./build-release-artifacts.sh
# into an official container.

set -e
set -x

CURDIR="$(dirname $BASH_SOURCE)"

. $CURDIR/_build-lib.sh

# Setup a temp docker build context dir
WKDIR=$(mktemp -d docker-build.XXXX)
cp $CURDIR/Dockerfile.runtime $WKDIR

BUILD_INFO=$(${CURDIR}/version-tool build-info)
VERSION_INFO=$(${CURDIR}/version-tool version)

# Hard code the platform dir here
cp $CURDIR/../_docker_workspace/out/$RELEASE_PLATFORM/bin/* $WKDIR/
mkdir -p $WKDIR/python
cp python/*.py $WKDIR/python/
cp python/cf-runtime-requirements.txt $WKDIR/
echo "{\"version\": \"${VERSION_INFO}\", \"build\": \"${BUILD_INFO}\"}" \
  > $WKDIR/python/VERSION_BUILD.json

echo "Docker build context:"
ls -la $WKDIR

VERSION_BUILD_ARGS=$(${CURDIR}/version-tool docker-build-args)
docker build --force-rm ${NO_CACHE_ARGS} \
  -t $IMG_TAG \
  --label BUILD_STAMP=$BUILD_STAMP \
  ${VERSION_BUILD_ARGS} \
  -f $WKDIR/Dockerfile.runtime \
  $WKDIR

docker history $IMG_TAG
docker inspect -f '{{ range $k, $v := .ContainerConfig.Labels -}}
{{ $k }}={{ $v }}
{{ end -}}' "$IMG_TAG"

echo "Built docker image $IMG_TAG"

rm -rf docker-build.????

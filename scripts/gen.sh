#!/usr/bin/env bash

export GOPATH="/Users/nevermore/go"

GENS="$1"

# The working directory which was the root path of our project.
ROOT_PACKAGE="github.com/nevercase/publisher"

if [ "${GENS}" = "api" ] || grep -qw "api" <<<"${GENS}"; then
  cp ${GOPATH}/bin/go-to-protobuf-api ${GOPATH}/bin/go-to-protobuf
  Packages="$ROOT_PACKAGE/pkg/types"
  "${GOPATH}/bin/go-to-protobuf" \
     --packages "${Packages}" \
     --clean=false \
     --only-idl=false \
     --keep-gogoproto=false \
     --verify-only=false \
     --proto-import ${GOPATH}/src/k8s.io/api/core/v1
fi
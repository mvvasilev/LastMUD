#!/bin/bash

set -e

pushd $(dirname "$0")/.. # run from root dir

go run cmd/lastmudserver/main.go

popd
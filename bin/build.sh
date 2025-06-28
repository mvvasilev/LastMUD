#!/bin/bash

set -e

pushd $(dirname "$0")/.. # run from root dir

go build -o target/lastmudserver cmd/lastmudserver/main.go 

popd
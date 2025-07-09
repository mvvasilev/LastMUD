#!/bin/bash

set -e

pushd $(dirname "$0")/.. # run from root dir

go build -o target/game_server cmd/game_server/main.go

popd
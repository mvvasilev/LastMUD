#!/bin/bash

set -e

pushd $(dirname "$0")/.. # run from root dir

go run cmd/game_server/main.go

popd
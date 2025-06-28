#!/bin/bash

pushd $(dirname "$0")/.. # run from root dir

rm -rf ./coverage
mkdir ./coverage/

go test --cover -coverpkg=./internal... -covermode=count -coverprofile=./coverage/cover.out ./...

go tool cover -html=./coverage/cover.out

popd # switch back to dir we started from
#!/bin/bash

set -e

pushd $(dirname "$0")/.. # run from root dir

docker build -t lastmud_game_server .

popd
#!/usr/bin/env bash

# run something in a directory
# goes hand-in .run-controller-gen.sh

cd "$1"
echo "DIR: $1"
shift
echo "CMD: $@"
"$@"

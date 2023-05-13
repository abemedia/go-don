#!/bin/bash

diff="$(git diff)"

if [ -n "$diff" ]; then
  git stash --keep-index --include-untracked && git stash drop
fi

task mod lint test
exitCode=$?

git add .

if [ -n "$diff" ]; then
  echo "$diff" | git apply
fi

exit $exitCode

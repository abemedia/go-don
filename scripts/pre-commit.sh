#!/bin/bash

diff="$(git diff)"

if [ -n "$diff" ]; then
  git stash --keep-index --include-untracked --quiet && git stash drop --quiet
fi

exitCode=0

go mod tidy || exitCode=$?
golangci-lint run || exitCode=$?
go test ./... || exitCode=$?

if [ $exitCode -eq 0 ]; then
  git add .
else
  git stash --keep-index --include-untracked --quiet && git stash drop --quiet
fi

if [ -n "$diff" ]; then
  echo "$diff" | git apply > /dev/null
fi

exit $exitCode

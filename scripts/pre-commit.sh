#!/bin/bash

git stash --keep-index --include-untracked --quiet

exitCode=0

go mod tidy || exitCode=$?
golangci-lint run || exitCode=$?
go test ./... || exitCode=$?

if [ $exitCode -eq 0 ]; then
  git add .
else
  git stash --keep-index --include-untracked --quiet && git stash drop --quiet
fi

git stash pop --quiet

exit $exitCode

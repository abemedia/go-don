#!/bin/bash

diff=$(git diff)

if [ -n "$diff" ]; then
  git stash --keep-index -u
fi

task
git add .

if [ -n "$diff" ]; then
  git stash pop
fi

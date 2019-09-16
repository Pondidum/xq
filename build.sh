#! /bin/bash

export GO111MODULE=on

git_commit=$(git rev-parse HEAD)
git_dirty=$(git status --porcelain)
git_dirty=${git_dirty:+"+CHANGES"}

ldflags="-X xq/version.GitCommit=$git_commit$git_dirty"

go get -v -d .
go build \
  -ldflags "$ldflags" \
  -o xq

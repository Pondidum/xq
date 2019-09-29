#! /bin/bash

export GO111MODULE=on

git_commit=$(git rev-parse HEAD)
git_dirty=$(git status --porcelain)
git_version=$(git tag -l --points-at HEAD | grep 'version/' | sed 's/version\/\(.*\)/\1/g' | sort | tail -n 1)

git_dirty=${git_dirty:+"+CHANGES"}

dev_tag=''

if [ "$git_version" == "" ]; then
  dev_tag="dev"
  git_version="0.0.0"
fi

ldflags="-X xq/version.GitCommit=$git_commit$git_dirty -X xq/version.VersionPrerelease=$dev_tag -X xq/version.Version=$git_version"

go get -v -d .
go build \
  -ldflags "$ldflags" \
  -o xq

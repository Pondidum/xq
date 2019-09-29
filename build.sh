#! /bin/bash

set -e

export GO111MODULE=on
echo "==> Read Environment Information"
echo "    Go Modules: $GO111MODULE"
echo "    CI_WINDOWS: $CI_WINDOWS"
echo "    CI_LINUX: $CI_LINUX"

echo "==> Read Git Infomation"

git_commit=$(git rev-parse HEAD)
git_dirty=$(git status --porcelain)
git_version=$(git tag -l --points-at HEAD | grep 'version/' | sed 's/version\/\(.*\)/\1/g' | sort | tail -n 1)

git_dirty=${git_dirty:+"+CHANGES"}

echo "    Commit: $git_commit"
echo "    Dirty: $git_dirty"
echo "    Version Tag: $git_version"

dev_tag=''

if [ "$git_version" == "" ]; then
  dev_tag="dev"
  git_version="0.0.0"
fi

ldflags="-X xq/version.GitCommit=$git_commit$git_dirty -X xq/version.VersionPrerelease=$dev_tag -X xq/version.Version=$git_version"

echo "==> Installing Dependencies"
go get -v -d .

echo "==> Building"
go build \
  -ldflags "$ldflags"

echo "==> Running Vet"

go vet

echo "==> Running Tests"

go test -v ./...

echo "==> Packaging Artifacts"

if [ "$CI_LINUX" = "true" ]; then
  echo "    Creating Linux Artifact"
  7z a "xq-$git_version-linux.zip" $APPVEYOR_BUILD_FOLDER/xq
elif [ "$CI_WINDOWS" = "True" ]; then
  echo "    Creating Windows Artifact"
  7z a "xq-$git_version-windows.zip" $APPVEYOR_BUILD_FOLDER/xq.exe
fi

echo "==> Done"
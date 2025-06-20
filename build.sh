#!/usr/bin/env bash

built_at="$(date +'%F %T %z')"
gitAuthor="vscodev"

if [[ -n "$1" && "$1" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
  version="${1#v}"
else
  version="dev"
fi

ldflags="\
-s -w \
-X 'github.com/vscodev/alist/v3/internal/conf.BuiltAt=$built_at' \
-X 'github.com/vscodev/alist/v3/internal/conf.GitAuthor=$gitAuthor' \
-X 'github.com/vscodev/alist/v3/internal/conf.Version=$version'"

build() {
  export CGO_ENABLED=1
  go build -ldflags="$ldflags"
}

build

#!/usr/bin/env bash

resolve_link() {
  "$(command -v greadlink || command -v readlink)" "$1"
}

abs_dirname() {
  local path="$1"
  local name
  local cwd
  cwd="$(pwd)"

  while [ -n "$path" ]; do
    cd "${path%/*}" || exit 1
    name="${path##*/}"
    path="$(resolve_link "$name" || true)"
  done

  pwd
  cd "$cwd" || exit 1
}

exec "$(dirname "$(abs_dirname "$0")")/Emacs" "$@"

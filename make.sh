#!/bin/env bash

dir=$(dirname "$0")
binary_name="ipt"
dist="dist"

build() {
  mkdir -p $dist
  go build -o "$dist/$binary_name" *.go
  chmod +x "$dist/$binary_name"
}

run() {
  go run *.go $@
}

help() {
  echo "Usage: $script_name <command>"
  echo "Commands:"
  declare -F | cut -d " " -f 3 | sed 's/^/    /'
}

cmd=${1:-help}
shift
$cmd "$@"

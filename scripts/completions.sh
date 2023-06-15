#!/bin/sh
set -e
export WARN_FLEEK=no
for sh in bash zsh fish; do
	echo "$sh"
	go run ./cmd/fleek/main.go completion "$sh" >"completions/fleek.$sh"
done

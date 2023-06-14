#!/bin/sh
set -e
rm -rf completions
export WARN_FLEEK=no
mkdir completions
for sh in bash zsh fish; do
	go run cmd/fleek/main.go completion "$sh" >"completions/fleek.$sh"
done
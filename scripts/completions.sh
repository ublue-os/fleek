#!/bin/sh
set -e
rm -rf completions
mkdir completions
for sh in bash zsh fish; do
	go run cmd/fleek/main.go completion "$sh" >"completions/goreleaser.$sh"
done
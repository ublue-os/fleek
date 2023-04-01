#!/bin/sh

rm -rf man/
for i in `find ./locales -type f`
do
    file=$(basename "$i" .yml)
    mkdir -p man/$file
    mkdir -p man/$file/man1
    LANG=$file go run ./cmd/fleek/main.go man | gzip -c -9 >man/$file/man1/fleek.1.gz
done



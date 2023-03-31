#!/bin/sh

rm -rf man/
for i in `find ./locales -type f`
do
    file=$(basename "$i" .yml)
    mkdir -p man/$file
    mkdir -p man/$file/man1
    LANG=$file $1 man > man/$file/man1/fleek.1 2> /dev/null
    gzip -6 man/$file/man1/fleek.1
done



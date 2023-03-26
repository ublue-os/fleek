set dotenv-load
CONTAINER_RUNNER := "podman"
CONTAINER_BUILDER := "buildah"


default: build



[private]
move +FILES:
  for FILE in {{FILES}}; do \
    mv "$FILE" "$FILE.bak" ; \
  done

[private]
unmove +FILES:
  for FILE in {{FILES}} ; do \
    mv "$FILE.bak" "$FILE" ; \
  done

[private]
cleanup +FILES:
  rm -rf {{FILES}}

backup: (move "$FLEEK_MANAGED/.fleek.yml" "$FLEEK_MANAGED/.config/home-manager")

clean: (cleanup "$FLEEK_MANAGED/.fleek.yml" "$FLEEK_MANAGED/.config/home-manager" "dist") 

restore: clean (unmove "$FLEEK_MANAGED/.fleek.yml" "$FLEEK_MANAGED/.config/home-manager")

default-env:
  cp .env.template .env

build:
  . .env
  go build -a -tags netgo -ldflags '-w -extldflags "-static"'

apply: 
  [ -e "./fleek" ] || just build
  ./fleek apply --push

man: 
  #!/bin/bash
  rm -rf man/
  for i in `find ./locales -type f`
  do
      file=$(basename "$i" .yml)
      echo "$file"
      just mkman "$file"
  done

mkman lang:
    mkdir -p man/{{lang}} ;  \
    mkdir -p man/{{lang}}/man1 ; \
    LANG={{lang}} ./fleek man > man/{{lang}}/man1/fleek.1 2> /dev/null ; \
    gzip -6 man/{{lang}}/man1/fleek.1 ; \

push: man (cleanup "fleek" "fleek.1" "fleek.1.gz")
  {{CONTAINER_BUILDER}} build --no-cache -t docker.io/bketelsen/fleek .
  {{CONTAINER_RUNNER}} push docker.io/bketelsen/fleek

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

deps:
  curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.52.2
  curl -sfL https://goreleaser.com/static/run | DISTRIBUTION=pro bash

lint:
  golangci-lint run

snapshot:
  goreleaser release --clean --snapshot

build:
  @source ./.env
  @go build -a -tags netgo -ldflags '-w -extldflags "-static"' github.com/ublue-os/fleek/cmd/fleek

apply: 
  [ -e "./fleek" ] || just build
  ./fleek apply --push

man: build
  ./man.sh ./fleek

push: man (cleanup "fleek" "fleek.1" "fleek.1.gz")
  {{CONTAINER_BUILDER}} build --no-cache -t docker.io/bketelsen/fleek .
  {{CONTAINER_RUNNER}} push docker.io/bketelsen/fleek

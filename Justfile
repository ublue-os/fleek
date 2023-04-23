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
example NAME:
  [ -e "./fleek" ] || just build
  @rm -rf examples/{{NAME}}
  @mkdir -p examples/{{NAME}}
  @./fleek generate --level {{NAME}} -l projects/ublue/fleek/examples/{{NAME}}

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

lint:
  golangci-lint run

build:
  @source ./.env
  @go build -a -tags netgo -ldflags '-w -extldflags "-static"' github.com/ublue-os/fleek/cmd/fleek

apply:
  [ -e "./fleek" ] || just build
  ./fleek apply --push

examples:
  [ -e "./fleek" ] || just build
  just example "high"
  just example "low"
  just example "default"
  just example "none"
  
completions:
  [ -e "./fleek" ] || just build
  mkdir -p completions
  ./fleek completion bash > completions/fleek.bash
  ./fleek completion zsh > completions/fleek.zsh
  ./fleek completion fish > completions/fleek.fish

man: build
  ./scripts/man.sh

tag version: lint build examples man completions
  ./scripts/create-release.sh {{version}}
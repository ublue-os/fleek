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
  mkdir -p grp
  cd grp
  wget -O grp.tar.gz https://github.com/goreleaser/goreleaser-pro/releases/download/v1.16.2-pro/goreleaser-pro_Linux_x86_64.tar.gz
  tar -xvf grp.tar.gz


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

examples:
  [ -e "./fleek" ] || just build
  rm -rf examples/none
  rm -rf examples/low
  rm -rf examples/default
  rm -rf examples/high
  mkdir -p examples/none
  mkdir -p examples/low
  mkdir -p examples/default
  mkdir -p examples/high
  ./fleek generate --level none -l projects/ublue/fleek/examples/none
  ./fleek generate --level default -l projects/ublue/fleek/examples/default
  ./fleek generate --level low -l projects/ublue/fleek/examples/low
  ./fleek generate --level high -l projects/ublue/fleek/examples/high
  
completions:
  [ -e "./fleek" ] || just build
  mkdir -p completions
  ./fleek completion bash > completions/fleek.bash
  ./fleek completion zsh > completions/fleek.zsh
  ./fleek completion fish > completions/fleek.fish

man: build
  ./scripts/man.sh

push: man (cleanup "fleek" "fleek.1" "fleek.1.gz")
  {{CONTAINER_BUILDER}} build --no-cache -t docker.io/bketelsen/fleek .
  {{CONTAINER_RUNNER}} push docker.io/bketelsen/fleek

set dotenv-load


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
  @./fleek generate --level {{NAME}} -l fleek/examples/{{NAME}}

[private]
unmove +FILES:
  for FILE in {{FILES}} ; do \
    mv "$FILE.bak" "$FILE" ; \
  done

[private]
cleanup +FILES:
  rm -rf {{FILES}}

clean: (cleanup "$FLEEK_MANAGED/.fleek.yml" "$FLEEK_MANAGED/.config/home-manager" "dist")

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
  @go build -a -tags netgo -ldflags '-w -extldflags "-static"' github.com/ublue-os/fleek/cmd/fleek

apply:
  [ -e "./fleek" ] || just build
  ./fleek apply --push

examples:
  [ -e "./fleek" ] || just build
  just example "none"
  just example "low"
  just example "default"
  just example "high"
  
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

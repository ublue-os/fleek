set dotenv-load


default: build

[private]
example NAME:
  [ -e "./fleek" ] || just build
  @rm -rf examples/{{NAME}}
  @mkdir -p examples/{{NAME}}
  @./fleek generate --level {{NAME}} -l projects/ublue/fleek/examples/{{NAME}}

default-env:
  cp .env.template .env

lint:
  golangci-lint run

snapshot:
  goreleaser release --clean --snapshot

build:
  @go build -a -tags netgo -ldflags '-w -extldflags "-static"' github.com/ublue-os/fleek/cmd/fleek

examples:
  [ -e "./fleek" ] || just build
  just example "none"
  just example "low"
  just example "default"
  just example "high"
  

tag version:
  ./scripts/create-release.sh {{version}}

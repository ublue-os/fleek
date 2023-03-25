set dotenv-load

default: build

[private]
move +FILES:
  for FILE in $FILES ; do \
    mv "$FILE" "$FILE.bak" ; \
  done

[private]
unmove +FILES:
  for FILE in $FILES ; do \
    mv "$FILE.bak" "$FILE" ; \
  done

[private]
cleanup +FILES:
  rm -rf $FILES

backup: (move "$FLEEK_MANAGED/.fleek.yml" "$FLEEK_MANAGED/.config/home-manager")

clean: (cleanup "$FLEEK_MANAGED/.fleek.yml" "$FLEEK_MANAGED/.config/home-manager" "dist") 

restore: clean (unmove "$FLEEK_MANAGED/.fleek.yml" "$FLEEK_MANAGED/.config/home-manager")

default-env:
  cp .env.template .env

build:
  go build -a -tags netgo -ldflags '-w -extldflags "-static"'

apply: 
  [ -e "./fleek" ] || just build
  ./fleek apply --push

push: (cleanup "fleek" "fleek.1" "fleek.1.gz")
  "$CONTAINER_BUILDER" build --no-cache -t docker.io/bketelsen/fleek .
  "$CONTAINER_RUNNER" push docker.io/bketelsen/fleek

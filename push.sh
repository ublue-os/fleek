#!/bin/bash
rm -f fleek
rm -f fleek.1
rm -f fleek.1.gz

podman build --no-cache -t docker.io/bketelsen/fleek .
podman push docker.io/bketelsen/fleek

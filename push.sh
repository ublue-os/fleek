#!/bin/bash
podman build --no-cache -t docker.io/bketelsen/fleek .
podman push docker.io/bketelsen/fleek

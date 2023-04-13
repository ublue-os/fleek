#!/usr/bin/env bash

# From https://github.com/Mic92/ssh-to-age/blob/main/bin/create-release.sh
# License : https://github.com/Mic92/ssh-to-age/blob/main/LICENSE

set -eu -o pipefail

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"
cd $SCRIPT_DIR/..

version=${1:-}
if [[ -z "$version" ]]; then
  echo "USAGE: $0 version" >&2
  exit 1
fi

if [[ "$(git symbolic-ref --short HEAD)" != "main" ]]; then
  echo "must be on main branch" >&2
  exit 1
fi

# ensure we are up-to-date
uncommited_changes=$(git diff --compact-summary)
if [[ -n "$uncommited_changes" ]]; then
  echo -e "There are uncommited changes, exiting:\n${uncommited_changes}" >&2
  exit 1
fi
git pull git@github.com:ublue-os/fleek main
unpushed_commits=$(git log --format=oneline origin/main..main)
if [[ "$unpushed_commits" != "" ]]; then
  echo -e "\nThere are unpushed changes, exiting:\n$unpushed_commits" >&2
  exit 1
fi
sed -i -e "s!version = \".*\"!version = \"${version}\"!" flake.nix
git add flake.nix
nix-build --no-out-link flake.nix
git commit -m "chore: bump version ${version}"
git tag -e "${version}"

echo 'now run `git push --tags origin main`'
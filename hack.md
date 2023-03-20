## Test build
nix run --impure home-manager/master -- -b bak build --flake .
## Apply Host
nix run --impure home-manager/master -- -b bak switch --flake .#bjk@ghanima

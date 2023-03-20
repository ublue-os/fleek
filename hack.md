nix run --impure home-manager/master -- -b bak build --flake .
nix run --impure home-manager/master -- -b bak switch --flake .#bjk@ghanima

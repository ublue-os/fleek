{ pkgs, misc, ... }: {
   home.shellAliases = {
    apply-ghanima = "nix run --impure home-manager/master -- -b bak switch --flake .#bjk@ghanima";
    
    fleeks = "cd /var/home/bjk/projects/ublue/fleek/examples/default";
    };
}

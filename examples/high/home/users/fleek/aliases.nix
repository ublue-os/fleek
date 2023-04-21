{ pkgs, misc, ... }: {
   home.shellAliases = {
    apply-beast = "nix run --impure home-manager/master -- -b bak switch --flake .#bjk@beast";
    
    fleeks = "cd ~/projects/ublue/fleek/examples/high";
    
    # bat --plain for unformatted cat
    catp = "bat -P";
    
    # replace cat with bat
    cat = "bat";
    };
}

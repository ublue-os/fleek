{ pkgs, misc, ... }: {
  home.packages = [
    # Fleek Bling
    pkgs.lazygit
    pkgs.jq
    pkgs.yq
    pkgs.neovim
    pkgs.neofetch
    pkgs.btop
    pkgs.cheat
  ];
}




{ pkgs, misc, ... }: {
  home.packages = [
    # Fleek Bling
    pkgs.fzf
    pkgs.ripgrep
    pkgs.vscode
  ];
}




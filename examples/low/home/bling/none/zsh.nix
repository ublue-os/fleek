{ pkgs, misc, ... }: {
  programs.zsh = {
    enableCompletion = true;
    enable = true;
    profileExtra = ''
    # Enable nix
    [ -r ~/.nix-profile/etc/profile.d/nix.sh ] && source  ~/.nix-profile/etc/profile.d/nix.sh
  '';
  };
}

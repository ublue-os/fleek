{
  programs.bash = {
    enable = true;
    enableCompletion = true;
    enableVteIntegration = true;
    profileExtra = ''
      # Nix
      [ -r ~/.nix-profile/etc/profile.d/nix.sh ] && source  ~/.nix-profile/etc/profile.d/nix.sh
    '';
  };
}

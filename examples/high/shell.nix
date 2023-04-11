{ pkgs, misc, ... }: {
    programs.exa.enableAliases = true;
    
    programs.exa.extraOptions = [
   "--group-directories-first"
   "--header"
];
    
    programs.bat.config = {
  theme = "TwoDark";
};
    # zsh
  programs.zsh.profileExtra = "[ -r ~/.nix-profile/etc/profile.d/nix.sh ] && source  ~/.nix-profile/etc/profile.d/nix.sh";
  programs.zsh.enableCompletion = true;
  programs.zsh.enable = true;
}

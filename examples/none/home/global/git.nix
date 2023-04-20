{ pkgs, lib, config, flake, ... }:
{
  programs.git = {
    enable = true;
    aliases = {
      pushall = "!git remote | xargs -L1 git push --all";
      graph = "log --decorate --oneline --graph";
      add-nowhitespace = "!git diff -U0 -w --no-color | git apply --cached --ignore-whitespace --unidiff-zero -";
    };
    userName = flake.config.people.users.${config.home.username}.name;
    userEmail = flake.config.people.users.${config.home.username}.email;
    extraConfig = {
      feature.manyFiles = true;
      init.defaultBranch = "main";
      gpg.format = "ssh";
    };

    signing = {
      key = flake.config.people.users.${config.home.username}.sshPrivateKey;
      signByDefault = builtins.stringLength flake.config.people.users.${config.home.username}.sshPrivateKey > 0;
    };

    lfs.enable = true;
    ignores = [ ".direnv" "result" ];
  };
}

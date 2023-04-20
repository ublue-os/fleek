{ inputs, lib, pkgs, config, outputs, flake, ... }:

{
  imports = [
    ./git.nix
  ];
  nixpkgs = {
    config = {
      allowUnfree = true;
      allowUnfreePredicate = (_: true);
    };
  };

  nix = {
    package = lib.mkDefault pkgs.nix;
    settings = {
      experimental-features = [ "nix-command" "flakes" "repl-flake" ];
      warn-dirty = false;
    };
  };

  systemd.user.startServices = "sd-switch";

  home.packages = with pkgs; [
    flake.inputs.fleek.packages.${pkgs.system}.default
  ];

  programs = {
    home-manager.enable = true;
    git.enable = true;
  };


}

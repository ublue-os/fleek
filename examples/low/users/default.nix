{ config, lib, ... }:
let
  userSubmodule = lib.types.submodule {
    options = {
      name = lib.mkOption {
        type = lib.types.str;
      };
      email = lib.mkOption {
        type = lib.types.str;
      };
      sshPublicKey = lib.mkOption {
        type = lib.types.str;
        description = ''
          SSH public key file path.
        '';
      };
      sshPrivateKey = lib.mkOption {
        type = lib.types.str;
        description = ''
          SSH private key file path.
        '';
      };
    outPath = lib.mkOption {
        type = lib.types.str;
      };
    };
  };
  peopleSubmodule = lib.types.submodule {
    options = {
      users = lib.mkOption {
        type = lib.types.attrsOf userSubmodule;
      };
    };
  };
in
{
  options = {
    people = lib.mkOption {
      type = peopleSubmodule;
    };
  };
  config = {
    people = import ./config.nix;
  };
}

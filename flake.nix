{
  description = "Fleek - 'Home as Code' for Humans";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/release-22.11";
  };

  outputs = { self, nixpkgs }: let
    # Current version
    version = "v0.9.3";
    # Supported systems
    systems = [
      "aarch64-linux" # 64-bit ARM Linux
      "x86_64-linux" # 64-bit Intel/AMD Linux
      "aarch64-darwin" # 64-bit ARM macOS
      "x86_64-darwin" # 64-bit Intel macOS
    ];
    # Helper for providing per-supported-system outputs
    forEachSystem = f: nixpkgs.lib.genAttrs systems (system: f {
      pkgs = import nixpkgs { inherit system; };
    });
  in {

    # Output Fleek as a Nix package so that others can easily install it using Nix:
    #
    # nix run github:ublue-os/fleek -- --help

    packages = forEachSystem ({ pkgs }: {
      default = pkgs.buildGoModule {
        pname = "fleek";
        inherit version;
        src = ./.;
        nativeBuildInputs = with pkgs; [
          installShellFiles # Shell completion helper function (see postInstall below)
        ];
        subPackages = [ "cmd/fleek" ];
        vendorSha256 = null;
        CGO_ENABLED = 0;
        ldflags = [
          "-s"
          "-w"
          "-X github.com/ublue-os/fleek/internal/build.Version=${version}"
          "-X github.com/ublue-os/fleek/internal/build.Commit=${self.rev}"
          "-X github.com/ublue-os/fleek/internal/build.CommitDate=1970-01-01T00:00:00Z"
          "-extldflags=-static"
        ];
        tags = [
          "netgo"
        ];
        postInstall = ''
          installShellCompletion --cmd fleek \
            --bash <($out/bin/fleek completion bash) \
            --fish <($out/bin/fleek completion fish) \
            --zsh <($out/bin/fleek completion zsh)
        '';
      };
    });
  };
}

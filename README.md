# fleek
`fleek`: the missing home-management tool for `nix`

## Getting Started

Install [Nix](https://nixos.org) then run:

```
fleek init
```

`fleek` reads your configuration from `$HOME~/.fleek.yml`. You can edit this
file with `fleek edit` which opens your `$EDITOR`.

Edit the list of packages you want installed manually by modifying `$HOME/.fleek.yml`.

## Usage

### Search

Search for nix packages with `fleek search <pattern>`, or at [nixos.org](https://search.nixos.org/packages).

### Apply Your Configuration

Modify $HOME/.fleek.yml to your liking and run `fleek apply`.

### Update Installed Packages

To update packages you've already installed run `fleek update`.
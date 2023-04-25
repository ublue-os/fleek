# Fleek - "Home as Code" for Humans

Fleek is an all-in-one management system for everything you need to be productive on your computer.

Status: BETA.  Probably won't eat your computer. Probably won't break your system, at least beyond simple recoverability.

What follows is a brief overview and quick start. See full documentation at [https://getfleek.dev](https://getfleek.dev).

Fleek v0.9.0 introduced some changes to how it should be installed. If you already have fleek installed, see [upgrade instructions](https://getfleek.dev/docs/upgrade) on the website.

![Fleek add](fleek-add.gif)

## Own your $HOME

### Instant Productivity

Fleek takes you from an empty slate to a fully productive working environment in less than five minutes.

### Take It With You

No matter whether you work on a shiny new M2 MacBook Air, a well-loved ThinkPad running Linux, or Windows with WSL, Fleek lets you take the exact same environment, tools, and configuration wherever you go.

### Zero Learning Curve To Start

You don't need to master a fancy DSL or spelunk through pages of online manuals to get started with Fleek. Answer two questions and you're instantly off to the races. Fleek gives you opinionated starter configurations for `bash` and `zsh` in four different levels of BLING. You can choose a standard close-to-stock experience, or dial your environment to 11 with all the latest desktop and terminal bling. And switching between them takes less than a minute when you change your mind.

### Every Tool At Your Fingertips

Whether you need to install a new programming language's toolset or the latest social media application, Fleek has you covered with the largest set of programs and packages in the world. Add a line to your `.fleek.yml` file and `fleek apply` yourself into freedom.

![Fleek search](fleek-search.gif)

### Eject Button Optional

If you reach a point where you've grown beyond Fleek's opinions and you want more, just `fleek eject` and manage your configurations manually.

## Party in the Front, Business in the Back

Fleek is a user-friendly wrapper around Nix and Nix Home Manager, but the friendly `fleek` command hides all the complexity from you. Edit a 10 line YAML file and Fleek harnesses the power of Nix behind the scenes.

## Getting Started

You need `nix`. We love the [Determinate Systems Installer](https://zero-to-nix.com/), but any `nix` is good. If you're on Fedora Silverblue [this script](https://github.com/dnkmmr69420/nix-installer-scripts/blob/main/installer-scripts/silverblue-nix-installer.sh) is the good stuff.

After Nix is installed you need to enable [flakes and the nix command](https://nixos.wiki/wiki/Flakes). It can be as simple as this:

```shell
mkdir -p ~/.config/nix
echo "experimental-features = nix-command flakes" >> ~/.config/nix/nix.conf
```
Rather than downloading and installing Fleek, we're going to use the power of `nix` to run it directly from GitHub. This will let Fleek manage itself, and you'll always have the latest version.

First let's bootstrap `fleek`:

```shell
nix run github:ublue-os/fleek -- init
```

This will create your configuration file and symlink it to `$HOME/.fleek.yml`. Open it with your favorite editor and take a look.

Make any changes to the `~/.fleek.yml` file you want... we recommend Bling Level `high` for the best experience.  

Now let's apply your configuration:

```shell
nix run github:ublue-os/fleek -- apply
```
It will take a bit to download and install everything, but when it's done you should see something like this:
       
```shell
... more text above this ...
Activating onFilesChange
Activating reloadSystemd
 [✓]  Operation completed successfully
```

*What happened here?* We just installed Nix Home Manager, configured it with your preferences, and applied it to your system. Note that up until this point you haven't even installed Fleek. The power of `nix` did it for you. 

You may need to close and re-open your terminal or even log out to see the changes.

## ~/.fleek.yml

Here's what mine looks like:

```file
───────┬───────────────────────────────────────────────────
       │ File: .fleek.yml
───────┼───────────────────────────────────────────────────
   1   │ aliases:
   2   │     fleeks: cd ~/.local/share/fleek
   3   │ bling: high
   4   │ ejected: false
   5   │ flakedir: .local/share/fleek
   6   │ git:
   7   │     autocommit: true
   8   │     autopull: true
   9   │     autopush: true
  10   │     enabled: true
  11   │ min_version: 0.8.4
  12   │ name: Fleek Configuration
  13   │ packages:
  14   │     - helix
  15   │     - go
  16   │     - gcc
  17   │     - nodejs
  18   │     - yarn
  19   │     - rustup
  20   │     - vhs
  21   │     - rnix-lsp
  22   │     - duf
  23   │ paths:
  24   │     - $HOME/bin
  25   │     - $HOME/.local/bin
  26   │     - $HOME/go/bin
  27   │ programs:
  28   │     - dircolors
  29   │ shell: zsh
  30   │ systems:
  31   │     - arch: aarch64
  32   │       hostname: f84d89911e5d.ant.amazon.com
  33   │       os: darwin
  34   │       username: brianjk
  35   │     - arch: x86_64
  36   │       hostname: beast
  37   │       os: linux
  38   │       username: bjk
  39   │ unfree: true
  40   │ users:
  41   │     - email: bketelsen@gmail.com
  42   │       name: Brian Ketelsen
  43   │       ssh_private_key_file: ~/.ssh/id_rsa
  44   │       ssh_public_key_file: ~/.ssh/id_rsa.pub
  45   │       username: brianjk
  46   │     - email: bketelsen@gmail.com
  47   │       name: Brian Ketelsen
  48   │       ssh_private_key_file: ~/.ssh/id_rsa
  49   │       ssh_public_key_file: ~/.ssh/id_rsa.pub
  50   │       username: bjk
───────┴──────────────────────────────────────────
```

I removed some of the aliases and systems just to make the example shorter.

Line 3: `bling: high` tells `fleek` how many extras I want in my $HOME setup. If you don't have a strong opinion I recommend `high`, because it isn't a really much stuff and the set we chose to add is really strong. Options are `none`, `low`, `default`, `high`.

Line 13: `packages:` starts a list of the packages I want installed. Mine are mostly focused around software development, but any package available in [nixpkgs](https://search.nixos.org/packages) is available. You can search for packages to install with the `fleek search` command.

Line 23: `paths:` starts a list of directories I want to add to my $PATH.

Line 27: `programs: ` - starts a list of programs to install. Programs are packages, but with optional configuration. See [the documentation](https://getfleek.dev/docs/programs) for more information.

Line 29: `shell: zsh` - this line isn't currently used, but will may be in the future. For now, it's just a placeholder.

Line 30: `systems:` These are added by `fleek` when you run `fleek init`, you shouldn't need to edit this part manually. Note that `fleek` and `nix` support macOS, Linux and WSL on Windows, so your configurations are fully portable.

Now that you've seen some of the possibile changes you can make, edit your `~/.fleek.yml` file and save it.

To apply your changes run `fleek apply`. `fleek` spins for a bit, and makes all the changes you requested. You may need to close and re-open your terminal application to see some of the changes, particularly if you add or remove fonts.


That's the quick start! From here, you can try `fleek add` to add packages from the CLI, `fleek search` to search for available packages. The full documentation is on the [fleek website](https://getfleek.dev).

### Behind the Scenes

Fancy animated gifs and long-winded README's are great, but what really happens when you run `fleek apply` the first time? I'm glad you asked...

1. `fleek` creates a [nix home-manager][def] configuration based on the [templates here](https://github.com/ublue-os/fleek/blob/main/internal/flake/templates/home.nix.tmpl).
1. `fleek` compiles the templates and writes them to disk at `~/.local/share/fleek` by default.
1. `fleek` calls the `nix` command, which does `nix` things to download and install all the packages in your `.fleek.yml` file.
1. The libraries and binaries you specify get installed in the `/nix` folder, and symlinked into your very own personal `nix` profile. This is stored in your $HOME directory under `~/.nix-profile`.
1. The `home-manager` command in the configuration assembles shell scripts and configurations for you based on the shell specified in your `.fleek.yml` file.
1. Your existing shell configuration files get renamed with a `.bak` extension, so you can go back to your business if `fleek` isn't for you.
1. New configurations are written, either `.zshrc` or `.bashrc`, and these configuration files add your `~/.nix-profile` and the packages you've installed to your $PATH

Fleek's goal is to give you 90% of the power of a fully customized, hand-written, artisinally crafted `home-manager` configuration -- without forcing you to learn `nix`, which can be (ok, it IS) intimidating.

You may like `fleek` and use it on one or more of your machines to install the things you need and use every day.

You might eventually be curious about how it all works. If that's the case you can always look at the nix flake that's generated for you and see what's going on behind the scenes. Cool! There's a lot you can do with `nix` that `fleek` doesn't even attempt.

Or you might be a practical pragmatist sort of computer user and not care about what's happening behind the curtain. That's fine with `fleek`. We're just here to help you go from zero to productive faster, more reliably, without a bunch of fuss.

## Motivation

`fleek` is born out of frustration.

It's annoying setting up a new computer just right. It's annoying having completely different configurations and tools each place I need to work. It's nearly impossible to configure a Mac and an Ubuntu server to have the same look & feel, much less the same tools.

It's a huge time sink spending time configuring everything, and even worse when you try to do it in a way you can share it between computers.

It wasn't until I was talking with Jorge Castro and he pointed out how complex my `home-manager` config was compared to a simple YAML file.

![nixmeup](nixmeup.png)

Jorge is a smart dude. That night the first bits of `fleek` were written.

Within 24 hours I took the plunge and dog-fooded `fleek` on my development laptop. That was a great moment. A day later I added `fleek remote` commands to push my configs to GitHub, then added the code to detect your arch & os and change the flake accordingly.

I spent the next two weeks tweaking the user experience and testing over and over. Countless virtual machines were provisioned and destroyed.

| *countless*

And now we're here. `fleek` is ready for a broader audience than me and Jorge - who is the biggest inspriation for all the features `fleek` has, and maybe just as importantly, doesn't have.

## Shoulders

Standing on the shoulders of giants:

This flake template was the thing that got everything started!

- [flake template](https://github.com/Misterio77/nix-starter-configs)
- [template license](https://github.com/Misterio77/nix-starter-configs/blob/main/LICENSE)

In my third rewrite, I looked at devbox and loved how they organized everything. I *borrowed* a LOT from this. And by *borrowed* I mean outright copy & pasted. Many supporting functions in this code were written by the JetPack team, and very lightly modified by me.

- [devbox](https://github.com/jetpack-io/devbox)
- [devbox license](https://github.com/jetpack-io/devbox/blob/main/LICENSE)

Tasty release shell script inspired by and copied from ssh-to-age.

- [ssh-to-age](https://github.com/Mic92/ssh-to-age/blob/main/bin/create-release.sh)
- [ssh-to-age license](https://github.com/Mic92/ssh-to-age/blob/main/LICENSE)

Fleek's development flake was created by [Luc Perkins](https://github.com/lucperkins).

Great terminal display code inspired by and copied from pcli.

- [pcli](https://github.com/pterm/pcli/)
- [pcli license](https://github.com/pterm/pcli/blob/main/LICENSE)

None of this is possible without Nix and Nix Home Manager:

- [nix](https://nixos.org/)
- [home manager][def]

[def]: https://github.com/nix-community/home-manager

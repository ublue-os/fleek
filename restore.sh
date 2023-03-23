#!/bin/sh
rm -f ~/.fleek.yml
rm -rf ~/.config/home-manager

mv ~/.fleek.yml.save ~/.fleek.yml
mv ~/.config/hm-save ~/.config/home-manager

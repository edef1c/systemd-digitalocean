{ buildGoPackage, fetchFromGitHub, lib, ... }:
(buildGoPackage rec {
  name = "systemd-digitalocean-generator";
  goPackagePath = name;
  src = builtins.filterSource (path: type: lib.hasSuffix ".go" path) ./.;
  goDeps = ./deps.nix;
}).bin

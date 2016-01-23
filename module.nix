{ config, pkgs, ... }:
let
  generator = pkgs.callPackage ./. {};
in {
  systemd.generators.systemd-digitalocean-generator = "${generator}/bin/systemd-digitalocean-generator";
}

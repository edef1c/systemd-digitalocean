{ config, pkgs, ... }:
let
  generator = pkgs.callPackage ./. {};
in {
  systemd.network.enable = true;
  systemd.generators.systemd-digitalocean-generator = "${generator}/bin/systemd-digitalocean-generator";
}

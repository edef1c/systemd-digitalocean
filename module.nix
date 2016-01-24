{ config, pkgs, ... }:
let
  generator = pkgs.callPackage ./. {};
in {
  networking.useDHCP = false;
  systemd.network.enable = true;
  systemd.services.digitalocean = {
    after = [ "network-pre.target" ];

    serviceConfig.Type = "notify";
    serviceConfig.ExecStart = "${generator}/bin/systemd-digitalocean-generator";

    serviceConfig.Restart = "on-failure";
    serviceConfig.RestartSec = 0;
  };
  systemd.services.systemd-networkd = {
    after = [ "digitalocean.service" ];
    wants = [ "digitalocean.service" ];
  };
}

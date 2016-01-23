{ goPackages, fetchFromGitHub, lib, ... }: with goPackages;
(buildGoPackage rec {
  name = "systemd-digitalocean-generator";
  goPackagePath = name;
  src = builtins.filterSource (path: type: lib.hasSuffix ".go" path) ./.;
  buildInputs = let
    netlink = buildGoPackage rec {
      owner  = "vishvananda";
      repo   = "netlink";
      rev    = "5a5eb317d73bc513ae0cfea7d9ac3c39f145e1db";
      sha256 = "171pv8cmw8iwm5ci7hdxrkxf0bj34hihr0q950jvilwyh8sasck6";

      name = repo;
      goPackagePath = "github.com/${owner}/${repo}";
      src = fetchFromGitHub { inherit owner repo rev sha256; };
    };
  in [ go-systemd netlink ];
}).bin

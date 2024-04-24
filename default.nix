{ pkgs ? (
    let
      inherit (builtins) fetchTree fromJSON readFile;
      inherit ((fromJSON (readFile ./flake.lock)).nodes) nixpkgs gomod2nix;
    in
    import (fetchTree nixpkgs.locked) {
      overlays = [
        (import "${fetchTree gomod2nix.locked}/overlay.nix")
      ];
    }
  )
, buildGoApplication ? pkgs.buildGoApplication
}:

buildGoApplication rec {
  pname = "junit-tools";
  version = "0.0.2";
  pwd = ./.;
  src = ./.;

  ldflags = [
    "-s"
    "-w"
    "-extldflags" "\"-fno-PIC -static\""
    "-X" "github.com/kadaan/junit-tools/version.Version=${version}"
    "-X" "github.com/kadaan/junit-tools/version.BuildUser=nix"
    "-X" "github.com/kadaan/junit-tools/version.BuildHost=localhost"
  ] ++ pkgs.lib.optionals pkgs.stdenv.isLinux [
    "-d"
  ];

  tags = [
    "osusergo"
  ] ++ pkgs.lib.optionals pkgs.stdenv.isLinux [
    "netgo"
    "static_build"
  ];

  modules = ./gomod2nix.toml;
}




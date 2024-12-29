let
  unstable = import (fetchTarball https://nixos.org/channels/nixos-unstable/nixexprs.tar.xz) { };
in
{ nixpkgs ? import <nixpkgs> {} }:
with nixpkgs; mkShell {
  buildInputs = with pkgs; [
    air
    go
    gopls
    nodejs
    foreman
#    liquidsoap
    postgresql_13
    python3
    ffmpeg
    wget
    inetutils
  ];
}

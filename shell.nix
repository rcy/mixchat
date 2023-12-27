{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  buildInputs = with pkgs; [
    air
    go
    gopls
    nodejs
    foreman
    liquidsoap
    postgresql_13
    python3
    ffmpeg
  ];
}

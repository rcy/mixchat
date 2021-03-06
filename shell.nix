{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  buildInputs = with pkgs; [
    nodejs
    foreman
    liquidsoap
    postgresql_13
    python3
    ffmpeg
  ];
}

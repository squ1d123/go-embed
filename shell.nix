{ pkgs ? import <nixpkgs> {} }:
  pkgs.mkShell {
      nativeBuildInputs = with pkgs.buildPackages; [ python312Packages.ansible-core ansible-lint];
}

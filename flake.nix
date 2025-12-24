{
  description = "A graphically stunning TUI system monitor in Go";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, utils }:
    utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
      in
      {
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            mage
            ttyd
            vhs
            gcc
          ];

          shellHook = ''
            echo "Rundown Dev Environment Loaded"
            echo "Available tools: go, ttyd, vhs"
          '';
        };
      }
    );
}

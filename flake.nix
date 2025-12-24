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
        packages.default = pkgs.buildGoModule {
          pname = "rundown";
          version = "0.1.0";
          src = ./.;
          vendorHash = "sha256-5CYrNtYvRjozfV6SQlm70EMFEqe7Ks84B5hI/EThxlI=";
        };

        apps.default = utils.lib.mkApp {
          drv = self.packages.${system}.default;
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            mage
            ttyd
            vhs
            gcc
            golangci-lint
          ];

          shellHook = ''
            echo "Rundown Dev Environment Loaded"
            echo "Available tools: go, mage, ttyd, vhs"
          '';
        };
      }
    );
}

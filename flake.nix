{
  description = "github.com/db757/iptools development environment";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs";
    flake-utils.url = "github:numtide/flake-utils";
    gomod2nix = {
      url = "github:nix-community/gomod2nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs = { self, nixpkgs, flake-utils, gomod2nix }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
        buildGoApplication =
          gomod2nix.legacyPackages.${system}.buildGoApplication;
        package_name = "iptools";
        bin_name = "ipt";
      in {
        packages.default = buildGoApplication {
          pname = "${package_name}";
          version = "0.1.0";
          src = ./.;
          pwd = ./.;

          vendorHash = null;
          postInstall = ''
            # Rename the binary
            mv $out/bin/${package_name} $out/bin/${bin_name}
          '';
        };

        devShells.default = pkgs.mkShell {
          name = "github.com/db757/${package_name}";
          hardeningDisable =
            [ "fortify" ]; # Fix for debugging go tests with CGO_ENABLED=1

          env = {
            PROJECT_NAME = "github.com/db757/${package_name}";
            # CGO_ENABLED = "0"; # Disable CGO for pure Go builds
          };

          buildInputs = with pkgs; [
            go
            gopls
            golangci-lint
            golangci-lint-langserver
            gotools
            govulncheck
            lefthook
            gomod2nix.packages.${system}.default # gomod2nix CLI
          ];

          shellHook = ''
            echo "$PROJECT_NAME dev shell is initialized."
            go version
          '';
        };
      });
}

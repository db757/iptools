{
  description = "github.com/db757/go-iprange development environment";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let pkgs = import nixpkgs { inherit system; };
      in {
        devShells.default = pkgs.mkShell {
          name = "github.com/db757/go-iprange";
          env = {
            PROJECT_NAME = "github.com/db757/go-iprange";
            # CGO_ENABLED = "0"; # Disable CGO for pure Go builds
          };

          buildInputs = with pkgs; [ go gopls ];

          shellHook = ''
            echo "$PROJECT_NAME dev shell is initialized."
            go version
          '';
        };
      });
}

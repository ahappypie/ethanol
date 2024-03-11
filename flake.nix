{
  description = "nix environment for ethanol, the unity catalog sql migration tool";

  # Flake inputs
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixpkgs-unstable";
  };

  # Flake outputs
  outputs = { self, nixpkgs }:
    let
      overlays = [
        (final: prev: {
          go = prev.go_1_22;
        })
        (final: prev: rec {
          databricks-cli = prev.callPackage ./nix/pkgs/databricks-cli {};
        })
       ];

      # Systems supported
      allSystems = [ "x86_64-linux" "aarch64-linux" "x86_64-darwin" "aarch64-darwin" ];

      # Helper to provide system-specific attributes
      forAllSystems = f: nixpkgs.lib.genAttrs allSystems (system: f {
        pkgs = import nixpkgs { inherit system overlays; };
      });
    in
    {
      # Development environment output
      devShells = forAllSystems ({ pkgs }: {
        default = pkgs.mkShell {
          # The Nix packages provided in the environment
          packages = with pkgs; [
            databricks-cli
            go
            git-filter-repo
          ];
        };
      });

      packages = forAllSystems ({ pkgs }: {
        default = pkgs.buildGo122Module rec {
          pname = "ethanol";
          version = "0.0.1";

          src = pkgs.lib.cleanSource self;

          vendorHash = "sha256-yWsBQ9qZDH5NkjnQyWqXDnnyF8zZZW1hMOlFC+fYFfc=";

          meta = {
            description = "A SQL Runner for Databricks Unity Catalog";
            homepage = "https://github.com/ahappypie/ethanol";
            maintainers = [ "ahappypie" ];
          };
        };
      });
    };
}

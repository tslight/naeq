{
  description = "The Secret Cipher of the UFOnauts as an API & CLI, because ¯\_(ツ)_/¯";

  inputs = {
    flake-schemas.url = "https://flakehub.com/f/DeterminateSystems/flake-schemas/*.tar.gz";
    nixpkgs.url = "https://flakehub.com/f/NixOS/nixpkgs/*.tar.gz";
  };

  # Flake outputs that other flakes can use
  outputs = { self, flake-schemas, nixpkgs }:
    let
      # Helpers for producing system-specific outputs
      supportedSystems = [ "x86_64-linux" "aarch64-darwin" "x86_64-darwin" "aarch64-linux" ];
      forEachSupportedSystem = f: nixpkgs.lib.genAttrs supportedSystems (system: f {
        pkgs = import nixpkgs { inherit system; };
      });
    in {
      # Schemas tell Nix about the structure of your flake's outputs
      schemas = flake-schemas.schemas;

      # Development environments
      devShells = forEachSupportedSystem ({ pkgs }: {
        default = pkgs.mkShell {
          # Pinned packages available in the environment
          packages = with pkgs; [
            gnumake
            go_1_20
            godef # jump to definition in editors
            golangci-lint # fast linter runners
            gotools # Go tools like goimports, godoc, and others
            gopls # go language server for using lsp plugins
            nixpkgs-fmt
          ];

          # A hook run every time you enter the environment
          shellHook = ''
            echo "Hello nix"
          '';
        };
      });
    };
}

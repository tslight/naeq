{ pkgs ? (
  let
    inherit (builtins) fetchTree fromJSON readFile;
    inherit ((fromJSON (readFile ./flake.lock)).nodes) nixpkgs gomod2nix;
  in
    import (fetchTree nixpkgs.locked) {
      overlays = [
        (import "${fetchTree gomod2nix.locked}/overlay.nix")
      ];
    }
)
, mkGoEnv ? pkgs.mkGoEnv
, gomod2nix ? pkgs.gomod2nix
}:

let
  goEnv = mkGoEnv { pwd = ./.; };
in
pkgs.mkShell {
  packages = with pkgs; [
    gnumake
    go_1_20
    godef # jump to definition in editors
    goEnv
    golangci-lint # fast linter runners
    gomod2nix
    gopls # go language server for using lsp plugins
    gotools # Go tools like goimports, godoc, and others
  ];
}

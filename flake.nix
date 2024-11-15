{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    { self
    , nixpkgs
    , flake-utils
    , ...
    } @ inputs:
    {
      overlays.default = final: prev: {
        confess = self.packages.${final.system}.default;
      };

      nixosModules.default = import ./module.nix;
    }
    //
    (flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs {
          inherit system;
        };
      in
      with pkgs; {
        devShells.default = mkShell {
          buildInputs = with pkgs; [
            go
            sqlite
          ];
        };

        packages.default = pkgs.callPackage ./build.nix { };
      }
    ));
}

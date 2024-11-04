{
  description = "Development environment flake";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
        {
          devShells.default = pkgs.mkShell {
            packages = with pkgs; [
              apple-sdk_11
              autoconf
              clang
              coreutils
              curl
              darwin.DarwinTools # sw_vers
              dbus
              expat
              findutils
              gcc
              giflib
              gmp
              gnumake
              gnused
              gnutar
              gnutls
              jansson
              lcms2
              libffi
              libgccjit
              libiconv
              libpng
              librsvg
              libtasn1
              libunistring
              libwebp
              libxml2
              mailutils
              ncurses
              nettle
              pkg-config
              python3
              rsync
              ruby_3_3
              sqlite
              texinfo
              time
              tree-sitter
              which
              xcbuild
              zlib
            ];

            shellHook = ''
            export MACOSX_DEPLOYMENT_TARGET="11.0"
            export DEVELOPER_DIR="${pkgs.apple-sdk_11}"
            export EMACS_BUILD_USE_NIX="true"
            export NIX_GCC_LIB_VERSION="${pkgs.gcc.cc.lib.version}"
            export NIX_GCC_LIB_ROOT="${pkgs.gcc.cc.lib.outPath}"
            export NIX_LIBGCCJIT_VERSION="${pkgs.libgccjit.version}"
            export NIX_LIBGCCJIT_ROOT="${pkgs.libgccjit.outPath}"
          '';
          };
        }
    );
}

{
  description = "Development environment flake";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/24.11-beta";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};

        # List of supported macOS SDK versions.
        sdk_versions = [ "11" "12" "13" "14" "15" ];
        default_sdk_version = "11";

        mkDevShell = { macos_version ? default_sdk_version }:
          let
            apple_sdk = pkgs.${"apple-sdk_${macos_version}"};
          in
            pkgs.mkShell {
              # Package list specifically excludes ncurses, so that we link
              # against the system version of ncurses. This ensures emacs' TUI
              # works out of the box without the user having to manually set
              # TERMINFO in the shell before launching emacs.
              packages = with pkgs; [
                apple_sdk
                autoconf
                bash
                cairo
                clang
                coreutils
                curl
                darwin.DarwinTools # sw_vers
                dbus
                expat
                findutils
                gcc
                gettext
                giflib
                git
                gmp
                gnumake
                gnupatch
                gnused
                gnutar
                gnutls
                harfbuzz
                jansson
                jq
                lcms2
                libffi
                libgccjit
                libiconv
                libjpeg
                libpng
                librsvg
                libtasn1
                libunistring
                libwebp
                libxml2
                mailutils
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
                export CC=clang
                export MACOSX_DEPLOYMENT_TARGET="${macos_version}.0"
                export DEVELOPER_DIR="${apple_sdk}"
                export NIX_LIBGCCJIT_VERSION="${pkgs.libgccjit.version}"
                export NIX_LIBGCCJIT_ROOT="${pkgs.libgccjit.outPath}"
                export BUNDLE_WITHOUT=development
              '';
            };

        # Generate an attrset of shells for each macOS SDK version.
        versionShells = builtins.listToAttrs (
          map (version: {
            name = "macos${version}";
            value = mkDevShell { macos_version = version; };
          }) sdk_versions
        );
      in
        {
          devShells = versionShells // {
            default = mkDevShell {};
          };
        }
    );
}

# build-emacs-for-macos

Use this script at your own risk. As of writing (2020-02-02) it works for me on
my own machine to build the `emacs-27` release branch. My machine is a late-2016
13-inch Touchbar MacBook Pro runnning macOS 10.15.2 and Xcode 11.3.

Your luck may vary.

The build produced does have some limitations:

- It is not a universal application. The CPU architecture of the built
  application will be that of the machine it was built on.
- The minimum required macOS version of the built application will be the same
  as that of the machine it was built on.

## Why?

- To use new features available from master or pre-release branches, which have
  not made it into a official stable release yet.
- Homebrew builds of Emacs are not self-contained applications, making it very
  difficult when doing HEAD builds and you need to rollback to a earlier
  version.
- Builds from [emacsformacosx.com](https://emacsformacosx.com/) has had no new
  nightly builds for two months right now.
- Both Homebrew HEAD builds, and nightly builds from emacsformacosx.com are
  built from the `master` branch. This script allows you to choose any branch
  you want. I am currently building from the `emacs-27` branch which is the
  basis of the upcoming Emacs 27 release, meaning it should be more stable than
  `master` builds.

## Requirements

- [Xcode](https://apps.apple.com/gb/app/xcode/id497799835?mt=12)
- [Homebrew](https://brew.sh/)
- All Homebrew formula listed in the `Brewfile`, which can all easily be
  installed by running:
  ```
  brew bundle
  ```

## Usage

Then to download a tarball of the `master` branch, build Emacs.app:

    ./build-emacs-for-macos

If you want to build the `emacs-27` git branch, run:

    ./build-emacs-for-macos emacs-27

If you want to build the stable `emacs-26.3` git tag, run:

    ./build-emacs-for-macos emacs-26.3

Resulting applications are saved to the `builds` directory in a bzip2 compressed
tarball.

## Internals

The script downloads the source code as a gzipped tar archive from the [GitHub
mirror](https://github.com/emacs-mirror/emacs) repository, as it makes it very
easy to get a tarball of any given git reference.

It then runs `./configure` with a various options, partly based on what [David
Caldwell](https://github.com/caldwell)'s
[build-emacs](https://github.com/caldwell/build-emacs) scripts do, including
copying various dynamic libraries into the application itself. So the built
application should in theory run on a macOS install that does not have homebrew,
or do no have the relevant brew formula installed.

Code quality, is well, non-existent. The build script started life a super-quick
hack back in 2013, and now it's even more of a dirty hack. I might clean it up
and add unit tests if I end up relying on this script for a prolonged period of
time. For now I plan to use it until Emacs 27 is officially released.

## License

```
        DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE
                    Version 2, December 2004

 Copyright (C) 2020 Jim Myhrberg

 Everyone is permitted to copy and distribute verbatim or modified
 copies of this license document, and changing it is allowed as long
 as the name is changed.

            DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE
   TERMS AND CONDITIONS FOR COPYING, DISTRIBUTION AND MODIFICATION

  0. You just DO WHAT THE FUCK YOU WANT TO.
```

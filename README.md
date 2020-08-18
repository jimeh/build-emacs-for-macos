# build-emacs-for-macos

My personal hacked together script for building a completely self-contained
Emacs.app application on macOS, from any git branch, tag, or ref.

Use this script at your own risk.

## Why?

- To use new features available from master or branches, which have not made it
  into a official stable release yet.
- Homebrew builds of Emacs are not self-contained applications, making it very
  difficult when doing HEAD builds and you need to rollback to a earlier
  version.
- Both Homebrew HEAD builds, and nightly builds from emacsformacosx.com are
  built from the `master` branch. This script allows you to choose any branch,
  tag, or git ref you want.

## Status

As of writing (2020-08-18) it works for me on my machine. Your luck may vary.

I have successfully built:

- `emacs-27.1` release git tag
- `master` branch (Emacs 28.x)
- `feature/native-comp` branch (Emacs 28.x)

For reference, my machine is:

- 13-inch MacBook Pro (2020)
- 10th Gen i7, 2.3 GHz 4-core/8-thread CPU
- macOS 10.15.6 (19G2021)
- Xcode 11.6

## Limitations

The build produced does have some limitations:

- It is not a universal application. The CPU architecture of the built
  application will be that of the machine it was built on.
- The minimum required macOS version of the built application will be the same
  as that of the machine it was built on.
- The application is not signed, so running it on machines other than the one
  that built the application will yield warnings. If you want to make a signed
  Emacs.app, google is you friend for finding signing instructions.

## Requirements

- [Xcode](https://apps.apple.com/gb/app/xcode/id497799835?mt=12)
- [Homebrew](https://brew.sh/)
- All Homebrew formula listed in the `Brewfile`, which can all easily be
  installed by running:
  ```
  brew bundle
  ```

## Usage

```
Usage: ./build-emacs-for-macos [options] <branch/tag/sha>

Branch, tag, and SHA are from the mirrors/emacs Github repo,
available here: https://github.com/mirrors/emacs
    -j, --parallel PROCS             Compile in parallel using PROCS processes
    -x, --xwidgets                   Apply XWidgets patch for Emacs 27
        --native-comp                Enable native-comp
        --native-fast-boot           Only relevant with --native-comp
```

Resulting applications are saved to the `builds` directory in a bzip2 compressed
tarball.

I would typically recommend to pass a `-j` value equal to the number of CPU
threads your machine has to ensure a fast build. In the below examples I'll be
using `-j 4`.

### Examples

To download a tarball of the `master` branch (Emacs 28.x), build Emacs.app from
it:

```
./build-emacs-for-macos -j 4
```

To build the stable `emacs-27.1` release git tag, with XWidgets support, run:

```
./build-emacs-for-macos -j 4 --xwidgets emacs-27.1
```

## Native-Comp

To build a Emacs.app with native-comp support
([gccemacs](https://akrl.sdf.org/gccemacs.html)) from the `feature/native-comp`
branch, you will need to install a patched version of Homebrew's `gcc` formula
that includes libgccjit.

The patch itself is in `./Formula/gcc.rb.patch`, and comes from
[this](https://gist.github.com/mikroskeem/0a5c909c1880408adf732ceba6d3f9ab#1-gcc-with-libgccjit-enabled)
gist.

You can install the patched formula by running the helper script:

```
./install-patched-gcc
```

The helper script will copy your local `gcc.rb` Forumla from Homebrew to
`./Formula`, and apply the `./Formula/gcc.rb.patch` to it. After which it then
proceed to install the patched gcc formula which includes libgccjit.

As it requires installing and compiling GCC from source, it can take anywhere
between 30-60 minutes or more depending on your machine.

And finally to build a Emacs.app with native compilation enabled, run:

```
./build-emacs-for-macos -j 4 --native-comp feature/native-comp
```

On my machine with `-j 8` this takes around 20-25 minutes. The increased build
time is cause all lisp files in the app are compiled to native `*.eln` files.

The build time can be sped up by using `--native-fast-boot`, which compiles a
minimal required set of lisp files to native code during build, and will compile
the rest dynamically in the background as they get loaded while you're using
Emacs.

## Credits

- I've borrowed some ideas and in general used
  [David Caldwell](https://github.com/caldwell)'s excellent
  [build-emacs](https://github.com/caldwell/build-emacs) project, which produces
  all builds for [emacsformacosx.com](https://emacsformacosx.com).
- Patches applied are pulled from
  [emacs-plus](https://github.com/d12frosted/homebrew-emacs-plus), which is an
  excellent Homebrew formula with lots of options not available elsewhere.
- The following gists were all extremely useful in figuring out how get get
  native-comp building on macOS:
  - https://gist.github.com/mikroskeem/0a5c909c1880408adf732ceba6d3f9ab#1-gcc-with-libgccjit-enabled
  - https://github.com/shshkn/emacs.d/blob/master/docs/nativecomp.md
  - https://gist.github.com/AllenDang/f019593e65572a8e0aefc96058a2d23e

## Internals

The script downloads the source code as a gzipped tar archive from the
[GitHub mirror](https://github.com/emacs-mirror/emacs) repository, as it makes
it very easy to get a tarball of any given git reference.

It then runs `./configure` with a various options, including copying various
dynamic libraries into the application itself. So the built application should
in theory run on a macOS install that does not have Homebrew, or does not have
the relevant Homebrew formulas installed.

Code quality of the script itself, is well, non-existent. The build script
started life a super-quick hack back in 2013, and now it's even more of a dirty
hack. I might clean it up and add unit tests if I end up relying on this script
for a prolonged period of time. For now I plan to use it at least until
native-comp lands in a stable Emacs release for macOS.

## License

[CC0 1.0 Universal](http://creativecommons.org/publicdomain/zero/1.0/)

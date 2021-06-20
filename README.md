# build-emacs-for-macos

My personal hacked together script for building a completely self-contained
Emacs.app application on macOS, from any git branch, tag, or ref. With support
for native-compilation.

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

As of writing (2021-04-25) it works for me on my machine. Your luck may vary.

I have successfully built:

- `emacs-27.1` release git tag
- `master` branch (Emacs 28.x)
- `feature/native-comp` branch (Emacs 28.x)

For reference, my machine is:

- 13-inch MacBook Pro (2020), 10th-gen 2.3 GHz Quad-Core Intel Core i7 (4c/8t)
- macOS Big Sur 11.2.3 (20D91)
- Xcode 12.4 (12D4e)

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
- Ruby 2.3.0 or later is needed to execute the build script itself. macOS comes
  with Ruby, check your version with `ruby --version`. If it's too old, you can
  install a newer version with:
  ```
  brew install ruby
  ```

## Usage

```
Usage: ./build-emacs-for-macos [options] <branch/tag/sha>

Branch, tag, and SHA are from the emacs-mirror/emacs/emacs Github repo,
available here: https://github.com/emacs-mirror/emacs

Options:
    -j, --parallel COUNT             Compile using COUNT parallel processes (detected: 8)
        --git-sha SHA                Override detected git SHA of specified branch allowing builds of old commits
        --[no-]xwidgets              Enable/disable XWidgets if supported (default: enabled)
        --[no-]native-comp           Enable/disable native-comp (default: enabled if supported)
        --[no-]native-march          Enable/disable -march=native CFLAG(default: disabled)
        --[no-]native-full-aot       Enable/disable NATIVE_FULL_AOT / Ahead of Time compilation (default: disabled)
        --[no-]rsvg                  Enable/disable SVG image support via librsvg (default: enabled)
        --no-titlebar                Apply no-titlebar patch (default: disabled)
        --no-frame-refocus           Apply no-frame-refocus patch (default: disabled)
        --[no-]github-auth           Make authenticated GitHub API requests if GITHUB_TOKEN environment variable is set.(default: enabled)
        --work-dir DIR               Specify a working directory where tarballs, sources, and builds will be stored and worked with
        --plan FILE                  Follow given plan file, instead of using given git ref/sha
```

Resulting applications are saved to the `builds` directory in a bzip2 compressed
tarball.

If you don't want the build process to eat all your CPU cores, pass in a `-j`
value of how many CPU cores you want it to use.

Re-building the same Git SHA again can yield weird results unless you first
trash the corresponding directory from the `sources` directory.

### Examples

To download a tarball of the `master` branch (Emacs 28.x with native-compilation
as of writing) and build Emacs.app from it:

```
./build-emacs-for-macos
```

To build the stable `emacs-27.1` release git tag run:

```
./build-emacs-for-macos emacs-27.1
```

All sources as downloaded as tarballs from the
[emacs-mirror](https://github.com/emacs-mirror/emacs) GitHub repository. Hence
to get a list of tags/branches available to install, simply check said
repository.

## Use Emacs.app as `emacs` CLI Tool

Builds come with a custom `emacs` shell script launcher for use from the command
line, located next to `emacsclient` in `Emacs.app/Contents/MacOS/bin`.

The custom `emacs` script makes sure to use the main
`Emacs.app/Contents/MacOS/Emacs` executable from the correct path, ensuring it
finds all the relevant dependencies within the Emacs.app bundle, regardless of
it it's exposed via `PATH` or symlinked to from elsewhere.

To use it, simply add `Emacs.app/Contents/MacOS/bin` to your `PATH`. For
example, if you place Emacs.app in `/Applications`:

```bash
if [ -d "/Applications/Emacs.app/Contents/MacOS/bin" ]; then
  export PATH="/Applications/Emacs.app/Contents/MacOS/bin:$PATH"
  alias emacs="emacs -nw" # Always launch "emacs" in terminal mode.
fi
```

If you want `emacs` in your terminal to launch a GUI instance of Emacs, don't
use the alias from the above example.

## Native-Comp

_Note: On 2021-04-25 the `feature/native-comp` branch was
[merged](http://git.savannah.gnu.org/cgit/emacs.git/commit/?id=289000eee729689b0cf362a21baa40ac7f9506f6)
into `master`._

The build script will automatically detect if the source tree being built
supports native-compilation, and enable it if available. You can override the
auto-detection logic to force enable or force disable native-compilation by
passing `--native-comp` or `--no-native-comp` respectfully.

By default `NATIVE_FULL_AOT` is disabled which ensures a fast build by native
compiling as few elisp source files as possible to build Emacs itself. Any
remaining elisp files will be dynamically compiled in the background the first
time they are used.

To enable native full Ahead-of-Time compilation, pass in the `--native-full-aot`
option, which will native-compile all of Emacs' elisp at built-time. On my
machine it takes around 10 minutes to build Emacs.app with `NATIVE_FULL_AOT`
disabled, and around 20-25 minutes with it enabled.

### Configuration

#### Native-Lisp Cache Directory

By default natively compiled `*.eln` files will be cached in
`~/.emacs.d/eln-cache/`. If you want to customize that, simply set a new path as
the first element of the `native-comp-eln-load-path` variable. The path string
must end with a `/`.

Below is an example which stores all compiled `*.eln` files in `cache/eln-cache`
within your Emacs configuration directory:

```elisp
(when (boundp 'native-comp-eln-load-path)
  (setcar native-comp-eln-load-path
          (expand-file-name "cache/eln-cache/" user-emacs-directory)))
```

#### Compilation Warnings

By default any warnings encountered during async native compilation will pop up
a warnings buffer. As this tends to happen rather frequently with a lot of
packages, it can get annoying. You can disable showing these warnings by setting
`native-comp-async-report-warnings-errors` to `nil`:

```elisp
(setq native-comp-async-report-warnings-errors nil)
```

### Issues

Please see all issues with the
[`native-comp`](https://github.com/jimeh/build-emacs-for-macos/issues?q=is%3Aissue+is%3Aopen+label%3Anative-comp)
label. It's a good idea if you read through them so you're familiar with the
types of issues and or behavior you can expect.

### Known Good Commits/Builds

A list of known "good" commits which produce working builds is tracked in:
[#6 Known good commits for native-comp](https://github.com/jimeh/build-emacs-for-macos/issues/6)

## Credits

- I've borrowed some ideas from [David Caldwell](https://github.com/caldwell)'s
  excellent [build-emacs](https://github.com/caldwell/build-emacs) project,
  which produces all builds for
  [emacsformacosx.com](https://emacsformacosx.com).
- Patches applied are pulled from
  [emacs-plus](https://github.com/d12frosted/homebrew-emacs-plus), which is an
  excellent Homebrew formula with lots of options not available elsewhere.
- The following sources were extremely useful in figuring out how get get the
  `feature/native-comp` branch building on macOS:
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

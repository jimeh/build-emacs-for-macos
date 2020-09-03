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

As of writing (2020-08-19) it works for me on my machine. Your luck may vary.

I have successfully built:

- `emacs-27.1` release git tag
- `master` branch (Emacs 28.x)
- `feature/native-comp` branch (Emacs 28.x)

For reference, my machine is:

- 13-inch MacBook Pro (2020), 10th-gen 2.3 GHz Quad-Core Intel Core i7 (4c/8t)
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
        --[no-]xwidgets              Enable/disable XWidgets (default: enabled)
        --[no-]native-comp           Enable/disable native-comp (default: enabled if supported)
        --[no-]native-fast-boot      Enable/disable NATIVE_FAST_BOOT (default: enabled if native-comp supported)
        --[no-]native-comp-macos-fixes
                                     Enable/disable fix based on feature/native-comp-macos-fixes branch (default: enabled if native-comp supported)
        --[no-]launcher              Enable/disable embedded launcher script  (default: enabled if native-comp is enabled)
```

Resulting applications are saved to the `builds` directory in a bzip2 compressed
tarball.

If you don't want the build process to eat all your CPU cores, pass in a `-j`
value of how many CPU cores you want it to use.

### Examples

To download a tarball of the `master` branch (Emacs 28.x as of writing) and
build Emacs.app from it:

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
./build-emacs-for-macos feature/native-comp
```

By default `NATIVE_FAST_BOOT` is enabled which ensures a fast build by native
compiling as few lisp source files as possible to build the app. Any remaining
lisp files will be dynamically compiled in the background the first time you use
them.

On my machine it takes around 10-15 minutes to build Emacs.app with
`NATIVE_FAST_BOOT` enabled. With it disabled it takes around 25 minutes.

### Configuration

Add the following near the top of your `early-init.el` or `init.el`:

```elisp
(setq comp-speed 2)
```

By default natively compiled `*.eln` files will be cached in
`~/.emacs.d/eln-cache/`. If you want to customize that, simply set a new path as
the first element of the `comp-eln-load-path` variable. The path string must end
with a `/`.

Also it seems somewhat common that some `*.eln` files are left behind with a
zero-byte file size if Emacs is quit while async native compilation is in
progress. Such empty files causes errors on startup, and needs to be deleted.

Below is an example which stores all compiled `*.eln` files in `cache/eln-cache`
within your Emacs configuration directory, and also deletes any `*.eln` files in
said directory which have a file size of zero bytes:

```elisp
(when (boundp 'comp-eln-load-path)
  (let ((eln-cache-dir (expand-file-name "cache/eln-cache/" user-emacs-directory))
        (find-exec (executable-find "find")))
    (setcar comp-eln-load-path eln-cache-dir)
    ;; Quitting emacs while native compilation in progress can leave zero byte
    ;; sized *.eln files behind. Hence delete such files during startup.
    (when find-exec
      (call-process find-exec nil nil nil eln-cache-dir
                    "-name" "*.eln" "-size" "0" "-delete"))))
```

### Issues

Please see all issues with the
[`native-comp`](https://github.com/jimeh/build-emacs-for-macos/issues?q=is%3Aissue+is%3Aopen+label%3Anative-comp)
label. It's a good idea if you read through them so you're familiar with the
types of issues and or behavior you can expect.

### Known Good Commits/Builds

A list of known "good" commits which produce working builds is tracked in:
[#6 Known good commits of feature/native-comp branch](https://github.com/jimeh/build-emacs-for-macos/issues/6)

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

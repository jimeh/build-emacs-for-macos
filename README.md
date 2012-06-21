# build-emacs-for-osx

Use this script at your own risk. It currently works for me on my own machine,
which as of writing runs:

* Mac OS X Lion 10.7.4 (11E53)
* Xcode + Command Line Tools 4.3.2 (4E2002)
* [GCC Installer][gcc] for 10.7

[gcc]: https://github.com/kennethreitz/osx-gcc-installer

Your luck might vary. Do note that it does not build a universal application.
The CPU architecture of the built application will be that of the machine it
was built on.

## Why?

I've been using [Homebrew](http://mxcl.github.com/homebrew/) the past few
months to build from HEAD. Homebrew comes with a number of patches, including
the [ns-toogle-fullscreen][fs] and [sRGB][] patches which I use.

Homebrew does not build a self-contained application though, which caused
issues for me when I needed to rollback to a specific build. I found the
easiest way to build a completely self-contained Emacs.app nightly from a
specific date with custom patches was to do it manually.

So I decided to quickly hack together a script to automate that manual
process. The code is a horrible hack, but it (seemingly) works as I'm writing
this in Emacs built with it.

## Usage

Myself I run the following command which will download a tarball of the
`master` branch, apply the fullscreen and sRGB patches, and build Emacs.app:

    ./build-emacs-for-osx

Or for example if you want to build the `EMACS_PRETEST_24_0_91` tag, run:

    ./build-emacs-for-osx EMACS_PRETEST_24_0_91

Resulting applications are saved to the `builds` directory in a bzip2
compressed tarball.

## Internals

I decided to pull Emacs' source from the GitHub mirror rather than the
official Bzr repo cause I'm not familiar with Bzr, and GitHub lets you easily
download tarballs of any branch, commit or tag.

The only option passed in `./configure` is `--with-ns`, meaning the resulting
application only supports the CPU architecture of the system is was built on.
There might be more side-effects to, but I haven't noticed any.


[fs]: https://gist.github.com/1012927
[srgb]: http://debbugs.gnu.org/cgi/bugreport.cgi?bug=8402

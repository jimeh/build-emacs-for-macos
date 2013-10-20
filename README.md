# build-emacs-for-osx

Use this script at your own risk. It currently works for me on my own machine,
which as of writing is:

* OS X 10.8.5 (12F45)
* Xcode 5.0 (5A1413)

Your luck might vary. Do note that it does not build a universal application.
The CPU architecture of the built application will be that of the machine it
was built on.


## Why?

I've been using [Homebrew][] the past few
months to build from HEAD. Homebrew comes with a number of patches, including
the [sRGB][] patches which I use.

[homebrew]: http://mxcl.github.com/homebrew/
[srgb]: http://debbugs.gnu.org/cgi/bugreport.cgi?bug=8402

Homebrew does not build a self-contained application though, which caused
issues for me when I needed to rollback to a specific build. I found the
easiest way to build a completely self-contained Emacs.app nightly from a
specific date with custom patches was to do it manually.

So I decided to quickly hack together a script to automate that manual
process. The code is a horrible hack, but it (seemingly) works as I'm writing
this in Emacs built with it.


## Usage

Myself I run the following command which will download a tarball of the
`master` branch, apply the sRGB patch, and build Emacs.app:

    ./build-emacs-for-osx

Or for example if you want to build the `emacs-24.3` tag, run:

    ./build-emacs-for-osx emacs-24.3

Resulting applications are saved to the `builds` directory in a bzip2
compressed tarball.


## Internals

I decided to pull Emacs' source from a GitHub [mirror][repo] rather than the
official Bzr repo cause I'm not familiar with Bzr, and GitHub lets you easily
download tarballs of any branch, commit or tag. And the tarballs from GitHub
are just over 30MB, compared to ~1GB to pull the offical Bzr repo.

[repo]: https://github.com/mirrors/emacs

The only option passed in `./configure` is `--with-ns`, meaning the resulting
application only supports the CPU architecture of the system is was built on.
There might be more side-effects to, but I haven't noticed any.


## License

```
        DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE
                    Version 2, December 2004

 Copyright (C) 2013 Jim Myhrberg

 Everyone is permitted to copy and distribute verbatim or modified
 copies of this license document, and changing it is allowed as long
 as the name is changed.

            DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE
   TERMS AND CONDITIONS FOR COPYING, DISTRIBUTION AND MODIFICATION

  0. You just DO WHAT THE FUCK YOU WANT TO.
```

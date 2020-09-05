# Changelog

All notable changes to this project will be documented in this file. See [standard-version](https://github.com/conventional-changelog/standard-version) for commit guidelines.

## 0.1.0 (2020-09-05)


### Features

* **deps:** add mailutils to Brewfile so Emacs can use GNU Mailutils ([d944a64](https://github.com/jimeh/build-emacs-for-macos/commit/d944a644847db564e56679b1c75c2550c8261b62))
* **native_comp:** add fix based on feature/native-comp-macos-fixes branch ([da2fcb0](https://github.com/jimeh/build-emacs-for-macos/commit/da2fcb0440a074a12c4fc6b1572cb55d8fb3cf9a))
* **native_comp:** Add support for --with-nativecomp ([fe460a8](https://github.com/jimeh/build-emacs-for-macos/commit/fe460a824ee57b602d29854167feab1b9f032aef))
* **native_comp:** embedd gcc/libgccjit into Emacs.app ([83289ac](https://github.com/jimeh/build-emacs-for-macos/commit/83289acd33b54a0d332fe92e2ad4ef7c1c633b72)), closes [#5](https://github.com/jimeh/build-emacs-for-macos/issues/5) [#7](https://github.com/jimeh/build-emacs-for-macos/issues/7)
* **native_comp:** support renaming of eln-cache director to native-lisp ([9d26435](https://github.com/jimeh/build-emacs-for-macos/commit/9d264357df61cf57a153947d2c22e28c27cba2d5))
* **patches:** add support for optional no-titlebar and no-refocus-frame patches ([583f22a](https://github.com/jimeh/build-emacs-for-macos/commit/583f22a360a08bf236ea0e0562e6fd1ddda3b933))
* **ref:** allow overriding git SHA ([eebda4d](https://github.com/jimeh/build-emacs-for-macos/commit/eebda4db42a700971b2083b3d420b99177f68b51))
* **release:** support building from release git tags ([c0e89b1](https://github.com/jimeh/build-emacs-for-macos/commit/c0e89b13648f9336aa46b2f088dcd439cb4028b7))


### Bug Fixes

* **deps:** Add missing dependencies to Brewfile ([39ea3eb](https://github.com/jimeh/build-emacs-for-macos/commit/39ea3eb5e8fed28f09962fcb8594ef85492b2f43))
* **native_comp:** ensure builds work after recent changes to eln cache locations ([b46e5aa](https://github.com/jimeh/build-emacs-for-macos/commit/b46e5aa7cba3c9496d2126fa9827275eaab720af)), closes [/akrl.sdf.org/gccemacs.html#org4b11ea1](https://github.com/jimeh//akrl.sdf.org/gccemacs.html/issues/org4b11ea1)
* **native_comp:** Improve ./install-patched-gcc helper ([a8d4db2](https://github.com/jimeh/build-emacs-for-macos/commit/a8d4db284cc216afa6793173f17e3811305eff05))
* **patches:** Fix patch download URL, add additional patches ([66acc01](https://github.com/jimeh/build-emacs-for-macos/commit/66acc01c0ca5d2f3c257f7df36082351a35f4273))
* **patches:** Only apply patches as part of archive extraction ([c4768f4](https://github.com/jimeh/build-emacs-for-macos/commit/c4768f4c3aed02a758863131abae5f07e8e4cf55))
* **requirements:** make script compatible with Ruby 2.3.0 and later ([8e459ce](https://github.com/jimeh/build-emacs-for-macos/commit/8e459ce00d8e5e3032ced260d8fbbc9b1dbc2c7a))
* **svg:** disable rsvg by default ([d30b45f](https://github.com/jimeh/build-emacs-for-macos/commit/d30b45fb2e507af98c3a958d159be3402a7a7bd1))
* **xwidgets:** Add support for emacs-27 specific xwidgets patch ([7767df0](https://github.com/jimeh/build-emacs-for-macos/commit/7767df0b660714c502d953b7bb22f5f3c2e3e3df))
* **xwidgets:** Use patch from emacs-plus Homebrew formula ([fb93beb](https://github.com/jimeh/build-emacs-for-macos/commit/fb93beb22c8b8dd0b46170c5c0b58159f25d6c1d))

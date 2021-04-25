# Changelog

All notable changes to this project will be documented in this file. See [standard-version](https://github.com/conventional-changelog/standard-version) for commit guidelines.

### [0.4.10](https://github.com/jimeh/build-emacs-for-macos/compare/0.4.9...0.4.10) (2021-04-25)


### Bug Fixes

* **cli:** correctly default to master branch if no git ref is given ([844df73](https://github.com/jimeh/build-emacs-for-macos/commit/844df73c8fa8440e657f7900ec89cdedb7c4c312))

### [0.4.9](https://github.com/jimeh/build-emacs-for-macos/compare/0.4.8...0.4.9) (2021-04-08)


### Bug Fixes

* **cli:** default to "master" if no git ref is given ([e19a6a7](https://github.com/jimeh/build-emacs-for-macos/commit/e19a6a7bc24379292ee06ae4c805b8c5365f2d97)), closes [#35](https://github.com/jimeh/build-emacs-for-macos/issues/35)
* **native_comp:** skip symlink creation for recent builds which do not need symlinks ([1000999](https://github.com/jimeh/build-emacs-for-macos/commit/1000999eb2673dc207a390ff3f902b9987b99173))

### [0.4.8](https://github.com/jimeh/build-emacs-for-macos/compare/0.4.7...0.4.8) (2021-02-27)


### Bug Fixes

* **native_comp:** add support for new --with-native-compilation flag ([581594d](https://github.com/jimeh/build-emacs-for-macos/commit/581594da3cfbf1dd2fa28e91710b767e21ff75d2))

### [0.4.7](https://github.com/jimeh/build-emacs-for-macos/compare/0.4.6...0.4.7) (2021-02-21)


### Bug Fixes

* **native_comp:** add libgccjit include dir during build stage ([e25ceaa](https://github.com/jimeh/build-emacs-for-macos/commit/e25ceaa7e25b0e1b9947401597845b5ba43e6cd1)), closes [#20](https://github.com/jimeh/build-emacs-for-macos/issues/20)

### [0.4.6](https://github.com/jimeh/build-emacs-for-macos/compare/0.4.5...0.4.6) (2021-02-15)


### Bug Fixes

* **native_comp:** improve env setup patch fixing potential issues ([dca023d](https://github.com/jimeh/build-emacs-for-macos/commit/dca023daecd8704f197cbc391380aa194bc47d62))

### [0.4.5](https://github.com/jimeh/build-emacs-for-macos/compare/0.4.4...0.4.5) (2021-01-06)


### Bug Fixes

* **cli:** remove defunct --[no-]native-comp-macos-fixes option ([ab55f54](https://github.com/jimeh/build-emacs-for-macos/commit/ab55f5421c81dc629e487bf4b8bb402657cb1bc4))

### [0.4.4](https://github.com/jimeh/build-emacs-for-macos/compare/0.4.3...0.4.4) (2021-01-02)


### Bug Fixes

* **deps:** add autoconf to Brewfile ([a47d3e0](https://github.com/jimeh/build-emacs-for-macos/commit/a47d3e0c6a8ea8161a3bad0eafdac2401cf53129))

### [0.4.3](https://github.com/jimeh/build-emacs-for-macos/compare/0.4.2...0.4.3) (2020-12-28)


### Bug Fixes

* **big-sur:** add Xcode CLI tools lib directory to runtime LIBRARY_PATH ([946856e](https://github.com/jimeh/build-emacs-for-macos/commit/946856e9c242d4a6fb5f839d8cae0acfafecdfc6))
* **big-sur:** added support for building on Big Sur ([2247158](https://github.com/jimeh/build-emacs-for-macos/commit/2247158051d0f59933569b6974b2b5269f13c79e))

### [0.4.2](https://github.com/jimeh/build-emacs-for-macos/compare/0.4.1...0.4.2) (2020-12-09)


### Bug Fixes

* **cli:** avoid error if --git-sha is used without a branch/tag/sha argument ([884f160](https://github.com/jimeh/build-emacs-for-macos/commit/884f1607f6707ca187b1abfb0ce562757d872230)), closes [#21](https://github.com/jimeh/build-emacs-for-macos/issues/21)
* **native_comp:** update env setup patch for recent changes to comp.el ([c7daa13](https://github.com/jimeh/build-emacs-for-macos/commit/c7daa1350bd69df172ce6484c54189d2cee8d97e))

### [0.4.1](https://github.com/jimeh/build-emacs-for-macos/compare/0.4.0...0.4.1) (2020-10-29)


### Features

* **native_comp:** remove patch based on feature/native-comp-macos-fixes branch ([70bf6b0](https://github.com/jimeh/build-emacs-for-macos/commit/70bf6b05d584976632b2fd2947c0bf692f5b6421))

## [0.4.0](https://github.com/jimeh/build-emacs-for-macos/compare/0.3.0...0.4.0) (2020-10-04)


### ⚠ BREAKING CHANGES

* **native_comp:** Standard Homewbrew `gcc` and `libgccjit` formula are now required for native-comp, instead of the custom patched gcc formula.

### Features

* **native_comp:** use new libgccjit Homebrew formula ([d8bbcb7](https://github.com/jimeh/build-emacs-for-macos/commit/d8bbcb72b33f6bde8678c9d37548217ffdf3d641))

## [0.3.0](https://github.com/jimeh/build-emacs-for-macos/compare/0.2.0...0.3.0) (2020-09-22)


### ⚠ BREAKING CHANGES

* **native_comp:** `--[no-]launcher` option is deprecated, as launcher script is no longer used.

### Features

* **native_comp:** use elisp patch instead of launcher script to set LIBRARY_PATH ([111cb64](https://github.com/jimeh/build-emacs-for-macos/commit/111cb6499368d14853a5927d38a43fc5e2f759f4)), closes [#14](https://github.com/jimeh/build-emacs-for-macos/issues/14)

## [0.2.0](https://github.com/jimeh/build-emacs-for-macos/compare/0.1.1...0.2.0) (2020-09-20)


### ⚠ BREAKING CHANGES

* **native_comp:** Deprecate `--[no-]native-fast-boot` option in favor of `--[no-]native-full-aot`

### Features

* **native_comp:** add support for NATIVE_FULL_AOT, replacing NATIVE_FAST_BOOT ([0ab94da](https://github.com/jimeh/build-emacs-for-macos/commit/0ab94da15309b04978982369bdfa17e03e9b6329))

### [0.1.1](https://github.com/jimeh/build-emacs-for-macos/compare/0.1.0...0.1.1) (2020-09-19)


### Bug Fixes

* **internal:** improve macOS version detection ([c89d0a0](https://github.com/jimeh/build-emacs-for-macos/commit/c89d0a0b73dfc82d918c326d89b141f8a2fc4de4)), closes [#13](https://github.com/jimeh/build-emacs-for-macos/issues/13)

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

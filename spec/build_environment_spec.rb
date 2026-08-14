# frozen_string_literal: true

require 'fileutils'
require 'tmpdir'

require_relative 'spec_helper'

unless defined?(Build)
  load File.expand_path('../build-emacs-for-macos', __dir__)
end

RSpec.describe 'Build environment' do
  around do |example|
    original_env = ENV.to_h
    begin
      example.run
    ensure
      ENV.replace(original_env)
    end
  end

  it 'uses a compatible system ncurses stub for Nix native builds' do
    with_build_environment(use_nix: true) do |build, root_dir, sdk_lib_dir|
      File.write(
        File.join(sdk_lib_dir, 'libncurses.tbd'),
        "targets: [ arm64e-macos, x86_64-macos ]\n"
      )

      library_path = build.send(:env_LIBRARY_PATH)
      stub_dir = File.join(root_dir, 'sdk-stubs')

      expect(library_path).to include(stub_dir)
      expect(library_path).not_to include(sdk_lib_dir)
      expect(File.read(File.join(stub_dir, 'libncurses.tbd'))).to eq(
        "targets: [ arm64-macos, x86_64-macos ]\n"
      )
    end
  end

  it 'includes the Command Line Tools library path for non-Nix builds' do
    with_build_environment(use_nix: false) do |build, _root_dir, sdk_lib_dir|
      expect(build.send(:env_LIBRARY_PATH)).to include(sdk_lib_dir)
    end
  end

  it 'reports a missing system ncurses stub for Nix builds' do
    with_build_environment(use_nix: true) do |build, _root_dir, sdk_lib_dir|
      expect { build.send(:env_LIBRARY_PATH) }.to raise_error(
        Error,
        "macOS system ncurses stub not found: #{sdk_lib_dir}/libncurses.tbd"
      )
    end
  end

  it 'uses Tree-sitter 0.25 pkg-config metadata for Nix Emacs 30 builds' do
    ENV['PKG_CONFIG_PATH'] = '/existing/pkgconfig'
    ENV['NIX_TREE_SITTER_025_ROOT'] = '/nix/tree-sitter-0.25'
    ENV['NIX_TREE_SITTER_026_ROOT'] = '/nix/tree-sitter-0.26'

    with_build_environment(use_nix: true, ref: 'emacs-30.2') do |build, *, **|
      expected = [
        '/nix/tree-sitter-0.25/lib/pkgconfig',
        '/existing/pkgconfig'
      ]
      expect(build.send(:env_PKG_CONFIG_PATH)).to eq(expected)
    end
  end

  it 'uses Tree-sitter 0.26 pkg-config metadata for Nix Emacs 31+ builds' do
    ENV['PKG_CONFIG_PATH'] = '/existing/pkgconfig'
    ENV['NIX_TREE_SITTER_025_ROOT'] = '/nix/tree-sitter-0.25'
    ENV['NIX_TREE_SITTER_026_ROOT'] = '/nix/tree-sitter-0.26'

    with_build_environment(use_nix: true, ref: 'emacs-32') do |build, *, **|
      expected = [
        '/nix/tree-sitter-0.26/lib/pkgconfig',
        '/existing/pkgconfig'
      ]
      expect(build.send(:env_PKG_CONFIG_PATH)).to eq(expected)
    end
  end

  it 'reports a missing selected Nix Tree-sitter root' do
    ENV.delete('NIX_TREE_SITTER_025_ROOT')
    ENV['NIX_TREE_SITTER_026_ROOT'] = '/nix/tree-sitter-0.26'

    with_build_environment(use_nix: true, ref: 'emacs-30.2') do |build, *, **|
      expect { build.send(:env_PKG_CONFIG_PATH) }.to raise_error(
        Error,
        'NIX_TREE_SITTER_025_ROOT is required for Emacs 30 Tree-sitter builds'
      )
    end
  end

  it 'selects versioned Homebrew Tree-sitter pkg-config metadata ' \
     'by Emacs major' do
    allow(OS).to receive(:version).and_raise('OS.version must not be used')

    with_build_environment(use_nix: false, ref: 'emacs-30.2') do |build, *, **|
      expect(build.send(:tree_sitter_pkg_config_path)).to eq(
        '/homebrew/opt/tree-sitter@0.25/lib/pkgconfig'
      )
    end

    with_build_environment(use_nix: false, ref: 'emacs-31') do |build, *, **|
      expect(build.send(:tree_sitter_pkg_config_path)).to eq(
        '/homebrew/opt/tree-sitter/lib/pkgconfig'
      )
    end
  end

  def with_build_environment(use_nix:, ref: 'emacs-30')
    Dir.mktmpdir('build-environment-test') do |root_dir|
      sdk_lib_dir = File.join(root_dir, 'host-sdk', 'usr', 'lib')
      FileUtils.mkdir_p(sdk_lib_dir)
      build = build_with_environment(
        use_nix: use_nix,
        ref: ref,
        root_dir: root_dir,
        sdk_lib_dir: sdk_lib_dir
      )

      yield build, root_dir, sdk_lib_dir
    end
  end

  def build_with_environment(use_nix:, ref:, root_dir:, sdk_lib_dir:)
    Build.allocate.tap do |build|
      build.instance_variable_set(:@root_dir, root_dir)
      build.instance_variable_set(:@ref, ref)
      build.instance_variable_set(
        :@options,
        { native_comp: true, tree_sitter: true, use_nix: use_nix }
      )
      build.instance_variable_set(
        :@gcc_info,
        Struct.new(:lib_dir, :darwin_lib_dir, :libgccjit_lib_dir).new(
          '/gcc/lib', '/gcc/lib/darwin', '/libgccjit/lib'
        )
      )
      build.define_singleton_method(:command_line_tools_library_path) do
        sdk_lib_dir
      end
      build.define_singleton_method(:host_arch) { 'arm64' }
      build.define_singleton_method(:brew_dir) { '/homebrew' }
    end
  end
end

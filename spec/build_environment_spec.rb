# frozen_string_literal: true

require 'fileutils'
require 'tmpdir'

require_relative 'spec_helper'

load File.expand_path('../build-emacs-for-macos', __dir__)

RSpec.describe 'Build environment' do
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

  def with_build_environment(use_nix:)
    Dir.mktmpdir('build-environment-test') do |root_dir|
      sdk_lib_dir = File.join(root_dir, 'host-sdk', 'usr', 'lib')
      FileUtils.mkdir_p(sdk_lib_dir)
      build = build_with_environment(
        use_nix: use_nix,
        root_dir: root_dir,
        sdk_lib_dir: sdk_lib_dir
      )

      yield build, root_dir, sdk_lib_dir
    end
  end

  def build_with_environment(use_nix:, root_dir:, sdk_lib_dir:)
    Build.allocate.tap do |build|
      build.instance_variable_set(:@root_dir, root_dir)
      build.instance_variable_set(
        :@options, { native_comp: true, use_nix: use_nix }
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
    end
  end
end

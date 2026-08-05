# frozen_string_literal: true

require 'minitest/autorun'
require 'fileutils'
require 'tmpdir'

load File.expand_path('../build-emacs-for-macos', __dir__)

class BuildEnvironmentTest < Minitest::Test
  GccInfo = Struct.new(:lib_dir, :darwin_lib_dir, :libgccjit_lib_dir)

  def test_nix_native_build_uses_a_compatible_system_ncurses_stub
    with_build_environment(use_nix: true) do |build, root_dir, sdk_lib_dir|
      File.write(
        File.join(sdk_lib_dir, 'libncurses.tbd'),
        "targets: [ arm64e-macos, x86_64-macos ]\n"
      )

      library_path = build.send(:env_LIBRARY_PATH)
      stub_dir = File.join(root_dir, 'sdk-stubs')

      assert_includes library_path, stub_dir
      refute_includes library_path, sdk_lib_dir
      assert_equal "targets: [ arm64-macos, x86_64-macos ]\n",
                   File.read(File.join(stub_dir, 'libncurses.tbd'))
    end
  end

  def test_non_nix_native_build_includes_command_line_tools_library_path
    with_build_environment(use_nix: false) do |build, _root_dir, sdk_lib_dir|
      assert_includes build.send(:env_LIBRARY_PATH), sdk_lib_dir
    end
  end

  def test_nix_build_reports_a_missing_system_ncurses_stub
    with_build_environment(use_nix: true) do |build, _root_dir, sdk_lib_dir|
      error = assert_raises(Error) { build.send(:env_LIBRARY_PATH) }

      assert_equal(
        "macOS system ncurses stub not found: #{sdk_lib_dir}/libncurses.tbd",
        error.message
      )
    end
  end

  private

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
        GccInfo.new('/gcc/lib', '/gcc/lib/darwin', '/libgccjit/lib')
      )
      build.define_singleton_method(:command_line_tools_library_path) do
        sdk_lib_dir
      end
      build.define_singleton_method(:host_arch) { 'arm64' }
    end
  end
end

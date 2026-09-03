# frozen_string_literal: true

require 'time'
require 'tmpdir'

require_relative 'spec_helper'

unless defined?(Build)
  load File.expand_path('../build-emacs-for-macos', __dir__)
end

RSpec.describe 'Emacs patch selection' do
  describe 'effective version detection' do
    it 'detects the version from the extracted source during builds' do
      Dir.mktmpdir('emacs-source') do |source_dir|
        File.write(
          File.join(source_dir, 'configure.ac'),
          "AC_INIT([GNU Emacs], [30.2], [bug-gnu-emacs@gnu.org], [],\n"
        )
        build = build_for(
          'emacs-32', version: nil, use_nix: true, source_dir: source_dir
        )
        allow(build).to receive(:github_api_get) do
          raise 'GitHub API should not be used when source is available'
        end
        allow(ENV).to receive(:fetch).and_call_original
        allow(ENV).to receive(:fetch)
          .with('NIX_TREE_SITTER_025_ROOT', nil)
          .and_return('/nix/tree-sitter-0.25')

        expect(build.send(:effective_version)).to eq(30)
        expect(build.send(:tree_sitter_pkg_config_path)).to eq(
          '/nix/tree-sitter-0.25/lib/pkgconfig'
        )
        patches = patch_basenames(build)
        expect(patches).to include('fix-macos-tahoe-scrolling.patch')
        expect(patches).not_to include('fix-ns-scroll-crash.patch')
      end
    end

    {
      'a known release ref' => ['emacs-30.2', {}],
      'a positional raw SHA' => ['deadbeef', {}],
      'a --git-sha override on master' => ['master', { git_sha: 'deadbeef' }]
    }.each do |description, (ref, options)|
      it "fetches the source version for #{description} without local source" do
        build = build_for(ref, version: nil, use_nix: true, **options)
        configure_ac = [
          "AC_INIT([GNU Emacs], [31.0.50], [bug-gnu-emacs@gnu.org], [],\n"
        ].pack('m0')
        allow(build).to receive(:github_api_get)
          .with(
            '/repos/emacs-mirror/emacs/contents/configure.ac?ref=test-sha'
          )
          .and_return({ content: configure_ac }.to_json)
        allow(ENV).to receive(:fetch).and_call_original
        allow(ENV).to receive(:fetch)
          .with('NIX_TREE_SITTER_027_ROOT', nil)
          .and_return('/nix/tree-sitter-0.27')

        expect(build.send(:effective_version)).to eq(31)
        expect(build.send(:tree_sitter_pkg_config_path)).to eq(
          '/nix/tree-sitter-0.27/lib/pkgconfig'
        )
        patches = patch_basenames(build)
        expect(patches).not_to include('fix-macos-tahoe-scrolling.patch')
        expect(patches).to include('fix-ns-scroll-crash.patch')
      end
    end

    it 'reports an unresolved source version' do
      build = build_for('deadbeef', version: nil)
      source = ['AC_INIT([Not Emacs], [1.0])'].pack('m0')
      allow(build).to receive(:github_api_get).and_return(
        { content: source }.to_json
      )

      expect { build.send(:effective_version) }.to raise_error(
        Error,
        'Failed to detect Emacs version from: test-sha'
      )
    end
  end

  describe 'default emacs-plus patches' do
    {
      'emacs-29.4' => %w[
        fix-window-role.patch
        system-appearance.patch
        round-undecorated-frame.patch
      ],
      'emacs-30.2' => %w[
        fix-window-role.patch
        system-appearance.patch
        round-undecorated-frame.patch
        fix-ns-x-colors.patch
        fix-macos-tahoe-scrolling.patch
      ],
      'emacs-31' => %w[
        system-appearance.patch
        round-undecorated-frame.patch
        fix-ns-x-colors.patch
        fix-ns-scroll-crash.patch
      ],
      'emacs-32' => %w[
        system-appearance.patch
        round-undecorated-frame.patch
        fix-ns-x-colors.patch
        fix-ns-scroll-crash.patch
      ]
    }.each do |ref, expected|
      it "selects the exact default patch set for #{ref}" do
        expect(patch_basenames(build_for(ref))).to eq(expected)
      end
    end

    it 'limits the macOS Tahoe scrolling patch to Emacs 30' do
      versions = (29..32).select do |version|
        patch_basenames(build_for("emacs-#{version}"))
          .include?('fix-macos-tahoe-scrolling.patch')
      end

      expect(versions).to eq([30])
    end

    it 'does not apply Tree-sitter compatibility source patches' do
      sources = (29..32).flat_map do |version|
        patch_sources(build_for("emacs-#{version}"))
      end

      expect(sources).not_to include(a_string_matching(/treesit/i))
    end
  end

  describe 'version-bounded optional patches' do
    it 'limits no-frame-refocus to Emacs 27 through 29' do
      selected_versions = (27..32).select do |version|
        patches = patch_basenames(
          build_for("emacs-#{version}", no_frame_refocus: true)
        )
        patches.include?('no-frame-refocus-cocoa.patch')
      end

      expect(selected_versions).to eq([27, 28, 29])
    end

    it 'retains fix-window-role for Emacs 31 commits before 2025-07-31' do
      build = build_for('emacs-31', date: Time.parse('2025-07-30 23:59:59Z'))

      expect(patch_basenames(build)).to include('fix-window-role.patch')
    end

    it 'omits fix-window-role from Emacs 31 on and after 2025-07-31' do
      build = build_for('emacs-31', date: Time.parse('2025-07-31 00:00:00Z'))

      expect(patch_basenames(build)).not_to include('fix-window-role.patch')
    end

    it 'does not attempt the alpha-background comparison patch on Emacs 32' do
      build = build_for('emacs-32', alpha_background: true)

      expect(patch_sources(build)).not_to include(a_string_matching(/alpha/))
    end

    it 'retains alpha-background support for Emacs 29 through 31' do
      selected_versions = (29..32).select do |version|
        patch_sources(build_for("emacs-#{version}", alpha_background: true))
          .any? { |url| url.include?('alpha') }
      end

      expect(selected_versions).to eq([29, 30, 31])
    end
  end

  def build_for(
    ref,
    date: Time.parse('2025-08-01'),
    source_dir: nil,
    version: ref[/emacs-(\d+)/, 1]&.to_i || (32 if ref == 'master'),
    **options
  )
    Build.allocate.tap do |build|
      build.instance_variable_set(:@ref, ref)
      build.instance_variable_set(:@source_dir, source_dir) if source_dir
      build.instance_variable_set(:@effective_version, version) if version
      build.instance_variable_set(
        :@options,
        {
          alpha_background: false,
          no_frame_refocus: false,
          no_titlebar: false,
          patches: [],
          xwidgets: false
        }.merge(options)
      )
      build.instance_variable_set(:@meta, { sha: 'test-sha', date: date })
    end
  end

  def patch_basenames(build)
    patch_sources(build).map { |source| File.basename(URI.parse(source).path) }
  end

  def patch_sources(build)
    build.send(:build_patches).map do |patch|
      patch[:url] || patch[:file]
    end.compact
  end
end

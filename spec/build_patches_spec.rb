# frozen_string_literal: true

require 'time'

require_relative 'spec_helper'

unless defined?(Build)
  load File.expand_path('../build-emacs-for-macos', __dir__)
end

RSpec.describe 'Emacs patch selection' do
  describe 'effective version detection' do
    {
      'emacs-29.4' => 29,
      'emacs-30' => 30,
      'emacs-31' => 31,
      'emacs-32' => 32,
      'master' => 32
    }.each do |ref, version|
      it "detects #{ref} as Emacs #{version}" do
        expect(build_for(ref).send(:effective_version)).to eq(version)
      end
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
      ],
      'emacs-32' => %w[
        system-appearance.patch
        round-undecorated-frame.patch
        fix-ns-x-colors.patch
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

  def build_for(ref, date: Time.parse('2025-08-01'), **options)
    Build.allocate.tap do |build|
      build.instance_variable_set(:@ref, ref)
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

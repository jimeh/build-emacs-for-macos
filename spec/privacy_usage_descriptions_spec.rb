# frozen_string_literal: true

require 'fileutils'
require 'json'
require 'open3'
require 'tmpdir'

require_relative 'spec_helper'
require_relative '../lib/privacy_usage_descriptions'

module PrivacyUsageDescriptionsSpecHelpers
  EXPECTED_DESCRIPTIONS = {
    'NSAppBundlesUsageDescription' =>
      'Emacs or a program running within it would like to access files ' \
      'inside other applications.',
    'NSAppDataUsageDescription' =>
      'Emacs or a program running within it would like to access files in ' \
      "other applications' data containers.",
    'NSAppleEventsUsageDescription' =>
      'A program running within Emacs would like to control other ' \
      'applications using Apple events.',
    'NSAppleMusicUsageDescription' =>
      'A program running within Emacs would like to access your media library.',
    'NSAudioCaptureUsageDescription' =>
      'A program running within Emacs would like to capture system audio.',
    'NSBluetoothAlwaysUsageDescription' =>
      'A program running within Emacs would like to use Bluetooth.',
    'NSCalendarsFullAccessUsageDescription' =>
      'A program running within Emacs would like to read and modify your ' \
      'calendars.',
    'NSCalendarsUsageDescription' =>
      'A program running within Emacs would like to access your calendars.',
    'NSCalendarsWriteOnlyAccessUsageDescription' =>
      'A program running within Emacs would like to add events to your ' \
      'calendars.',
    'NSCameraUsageDescription' =>
      'A program running within Emacs would like to use the camera.',
    'NSContactsUsageDescription' =>
      'A program running within Emacs would like to access your contacts.',
    'NSDesktopFolderUsageDescription' =>
      'Emacs or a program running within it would like to access files in ' \
      'Desktop folder.',
    'NSDocumentsFolderUsageDescription' =>
      'Emacs or a program running within it would like to access files in ' \
      'Documents folder.',
    'NSDownloadsFolderUsageDescription' =>
      'Emacs or a program running within it would like to access files in ' \
      'Downloads folder.',
    'NSFileProviderDomainUsageDescription' =>
      'Emacs or a program running within it would like to access files ' \
      'managed by file providers.',
    'NSLocalNetworkUsageDescription' =>
      'A program running within Emacs would like to access devices on your ' \
      'local network.',
    'NSLocationAlwaysAndWhenInUseUsageDescription' =>
      'A program running within Emacs would like to access your location ' \
      'even when Emacs is not active.',
    'NSLocationUsageDescription' =>
      'A program running within Emacs would like to access your location.',
    'NSLocationWhenInUseUsageDescription' =>
      'A program running within Emacs would like to access your location ' \
      'while Emacs is in use.',
    'NSMicrophoneUsageDescription' =>
      'A program running within Emacs would like to use your microphone.',
    'NSMotionUsageDescription' =>
      'A program running within Emacs would like to access motion data.',
    'NSNetworkVolumesUsageDescription' =>
      'Emacs or a program running within it would like to access files on ' \
      'volumes.',
    'NSPhotoLibraryAddUsageDescription' =>
      'A program running within Emacs would like to add items to your photo ' \
      'library.',
    'NSPhotoLibraryUsageDescription' =>
      'A program running within Emacs would like to access your photo library.',
    'NSRemindersFullAccessUsageDescription' =>
      'A program running within Emacs would like to read and modify your ' \
      'reminders.',
    'NSRemindersUsageDescription' =>
      'A program running within Emacs would like to access your reminders.',
    'NSRemovableVolumesUsageDescription' =>
      'Emacs or a program running within it would like to access files on ' \
      'removable volumes.',
    'NSSpeechRecognitionUsageDescription' =>
      'A program running within Emacs would like to use speech recognition.',
    'NSSystemAdministrationUsageDescription' =>
      'A program running within Emacs would like to modify system ' \
      'configuration.'
  }.freeze

  class MemoryEditor
    attr_reader :values

    def initialize(values = {})
      @values = values
    end

    def replace_string(_filename, key, value)
      values[key] = value
    end
  end
end

RSpec.describe PrivacyUsageDescriptions do
  it 'defines the expected privacy descriptions' do
    expect(described_class::DESCRIPTIONS).to eq(
      PrivacyUsageDescriptionsSpecHelpers::EXPECTED_DESCRIPTIONS
    )
  end

  it 'embeds descriptions without removing unrelated properties' do
    with_info_plist do |app, _info_plist|
      editor = PrivacyUsageDescriptionsSpecHelpers::MemoryEditor.new(
        'CFBundleName' => 'Emacs',
        'NSCameraUsageDescription' => 'Upstream description'
      )

      described_class.new(app, editor: editor).embed

      expect(editor.values['CFBundleName']).to eq('Emacs')
      expect(editor.values).to eq(
        PrivacyUsageDescriptionsSpecHelpers::EXPECTED_DESCRIPTIONS.merge(
          'CFBundleName' => 'Emacs'
        )
      )
    end
  end

  it 'rejects an app without an Info.plist' do
    Dir.mktmpdir('privacy-descriptions-test') do |dir|
      expect do
        described_class.new(
          File.join(dir, 'Emacs.app'),
          editor: PrivacyUsageDescriptionsSpecHelpers::MemoryEditor.new
        ).embed
      end.to raise_error(described_class::Error, /Info\.plist not found/)
    end
  end

  it 'replaces string values with plutil' do
    command = nil
    editor = described_class::PlutilEditor.new do |*args|
      command = args
    end

    editor.replace_string(
      '/tmp/Emacs.app/Contents/Info.plist',
      'NSCameraUsageDescription',
      'Camera reason'
    )

    expect(command).to eq([
                            'plutil', '-replace',
                            'NSCameraUsageDescription', '-string',
                            'Camera reason',
                            '/tmp/Emacs.app/Contents/Info.plist'
                          ])
  end

  it 'produces the expected manifest with the real plutil' do
    skip 'macOS plutil integration test' unless RUBY_PLATFORM.include?('darwin')

    with_info_plist do |app, info_plist|
      runner = lambda do |*args|
        _stdout, stderr, status = Open3.capture3(*args)
        raise stderr unless status.success?
      end
      editor = described_class::PlutilEditor.new(&runner)

      described_class.new(app, editor: editor).embed

      stdout, stderr, status =
        Open3.capture3('plutil', '-convert', 'json', '-o', '-', info_plist)
      expect(status).to be_success, stderr

      values = JSON.parse(stdout)
      expect(values.delete('CFBundleName')).to eq('Emacs')
      expect(values).to eq(
        PrivacyUsageDescriptionsSpecHelpers::EXPECTED_DESCRIPTIONS
      )
      expect(values.values).to all(be_a(String))
    end
  end

  def with_info_plist
    Dir.mktmpdir('privacy-descriptions-test') do |dir|
      app = File.join(dir, 'Emacs.app')
      contents = File.join(app, 'Contents')
      info_plist = File.join(contents, 'Info.plist')
      FileUtils.mkdir_p(contents)
      File.write(info_plist, <<~PLIST)
        <?xml version="1.0" encoding="UTF-8"?>
        <!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" \
          "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
        <plist version="1.0">
        <dict>
          <key>CFBundleName</key>
          <string>Emacs</string>
        </dict>
        </plist>
      PLIST

      yield app, info_plist
    end
  end
end

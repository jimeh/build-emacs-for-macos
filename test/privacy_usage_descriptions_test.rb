# frozen_string_literal: true

require 'json'
require 'minitest/autorun'
require 'open3'
require 'fileutils'
require 'tmpdir'

require_relative '../lib/privacy_usage_descriptions'

class PrivacyUsageDescriptionsTest < Minitest::Test
  EXPECTED_DESCRIPTIONS = {
    'NSAppBundlesUsageDescription' =>
      'A program running within Emacs would like to access files inside ' \
      'other applications.',
    'NSAppDataUsageDescription' =>
      'A program running within Emacs would like to access files in other ' \
      "applications' data containers.",
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
      'A program running within Emacs would like to access files in your ' \
      'Desktop folder.',
    'NSDocumentsFolderUsageDescription' =>
      'A program running within Emacs would like to access files in your ' \
      'Documents folder.',
    'NSDownloadsFolderUsageDescription' =>
      'A program running within Emacs would like to access files in your ' \
      'Downloads folder.',
    'NSFileProviderDomainUsageDescription' =>
      'A program running within Emacs would like to access files managed by ' \
      'file providers.',
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
      'A program running within Emacs would like to access files on network ' \
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
      'A program running within Emacs would like to access files on ' \
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

  def test_defines_the_expected_privacy_descriptions
    assert_equal EXPECTED_DESCRIPTIONS, PrivacyUsageDescriptions::DESCRIPTIONS
  end

  def test_embeds_descriptions_without_removing_unrelated_properties
    with_info_plist do |app, _info_plist|
      editor = MemoryEditor.new(
        'CFBundleName' => 'Emacs',
        'NSCameraUsageDescription' => 'Upstream description'
      )

      PrivacyUsageDescriptions.new(app, editor: editor).embed

      assert_equal 'Emacs', editor.values['CFBundleName']
      assert_equal EXPECTED_DESCRIPTIONS.merge('CFBundleName' => 'Emacs'),
                   editor.values
    end
  end

  def test_rejects_an_app_without_an_info_plist
    Dir.mktmpdir('privacy-descriptions-test') do |dir|
      error = assert_raises(PrivacyUsageDescriptions::Error) do
        PrivacyUsageDescriptions.new(File.join(dir, 'Emacs.app'),
                                     editor: MemoryEditor.new).embed
      end

      assert_match(/Info\.plist not found/, error.message)
    end
  end

  def test_plutil_editor_replaces_string_values
    command = nil
    editor = PrivacyUsageDescriptions::PlutilEditor.new do |*args|
      command = args
    end

    editor.replace_string('/tmp/Emacs.app/Contents/Info.plist',
                          'NSCameraUsageDescription', 'Camera reason')

    assert_equal [
      'plutil', '-replace', 'NSCameraUsageDescription', '-string',
      'Camera reason', '/tmp/Emacs.app/Contents/Info.plist'
    ], command
  end

  def test_real_plutil_produces_the_expected_manifest
    skip 'macOS plutil integration test' unless RUBY_PLATFORM.include?('darwin')

    with_info_plist do |app, info_plist|
      runner = lambda do |*args|
        _stdout, stderr, status = Open3.capture3(*args)
        raise stderr unless status.success?
      end
      editor = PrivacyUsageDescriptions::PlutilEditor.new(&runner)

      PrivacyUsageDescriptions.new(app, editor: editor).embed

      stdout, stderr, status =
        Open3.capture3('plutil', '-convert', 'json', '-o', '-', info_plist)
      assert status.success?, stderr

      values = JSON.parse(stdout)
      assert_equal 'Emacs', values.delete('CFBundleName')
      assert_equal EXPECTED_DESCRIPTIONS, values
      assert(values.values.all? { |value| value.is_a?(String) })
    end
  end

  private

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

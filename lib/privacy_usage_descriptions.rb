# frozen_string_literal: true

# Adds privacy usage descriptions needed when Emacs acts as the responsible
# application for Emacs Lisp or child processes that access protected macOS
# resources. Unknown Info.plist keys are harmless on older macOS versions, so
# every build receives the complete compatibility set.
class PrivacyUsageDescriptions
  class Error < StandardError
  end

  # These descriptions intentionally identify the program running within
  # Emacs. Emacs cannot know the program's specific reason for requesting a
  # resource, which is the same constraint faced by terminal applications.
  DESCRIPTIONS = {
    # Access to other application bundles, available on macOS 13 and later.
    'NSAppBundlesUsageDescription' =>
      'A program running within Emacs would like to access files inside ' \
      'other applications.',

    # Access to other applications' sandbox containers, macOS 14 and later.
    'NSAppDataUsageDescription' =>
      'A program running within Emacs would like to access files in other ' \
      "applications' data containers.",

    # Sending Apple events for automation, required since macOS 10.14.
    'NSAppleEventsUsageDescription' =>
      'A program running within Emacs would like to control other ' \
      'applications using Apple events.',

    # Access to the media library, required on macOS 15 and later for programs
    # linked against the macOS 15 SDK or later.
    'NSAppleMusicUsageDescription' =>
      'A program running within Emacs would like to access your media library.',

    # Capture of system audio, available on macOS 14.2 and later.
    'NSAudioCaptureUsageDescription' =>
      'A program running within Emacs would like to capture system audio.',

    # Bluetooth access, available on macOS 11 and later. The older
    # NSBluetoothPeripheralUsageDescription key is for iOS, not macOS.
    'NSBluetoothAlwaysUsageDescription' =>
      'A program running within Emacs would like to use Bluetooth.',

    # Full calendar access on macOS 14 and later.
    'NSCalendarsFullAccessUsageDescription' =>
      'A program running within Emacs would like to read and modify your ' \
      'calendars.',

    # Legacy calendar access for macOS 10.14 through macOS 13. Keep this with
    # the macOS 14 replacement keys to support the full deployment range.
    'NSCalendarsUsageDescription' =>
      'A program running within Emacs would like to access your calendars.',

    # Write-only calendar access on macOS 14 and later.
    'NSCalendarsWriteOnlyAccessUsageDescription' =>
      'A program running within Emacs would like to add events to your ' \
      'calendars.',

    # Camera access, required since macOS 10.14.
    'NSCameraUsageDescription' =>
      'A program running within Emacs would like to use the camera.',

    # Contacts access, required since macOS 10.8.
    'NSContactsUsageDescription' =>
      'A program running within Emacs would like to access your contacts.',

    # Protected Desktop access, available since macOS 10.15. The description
    # is optional but gives the system a useful prompt instead of a default.
    'NSDesktopFolderUsageDescription' =>
      'A program running within Emacs would like to access files in your ' \
      'Desktop folder.',

    # Protected Documents access, available since macOS 10.15.
    'NSDocumentsFolderUsageDescription' =>
      'A program running within Emacs would like to access files in your ' \
      'Documents folder.',

    # Protected Downloads access, available since macOS 10.15.
    'NSDownloadsFolderUsageDescription' =>
      'A program running within Emacs would like to access files in your ' \
      'Downloads folder.',

    # Access to files managed by a File Provider, macOS 10.15 and later.
    'NSFileProviderDomainUsageDescription' =>
      'A program running within Emacs would like to access files managed by ' \
      'file providers.',

    # Local-network privacy prompts, available on macOS 11 and later.
    'NSLocalNetworkUsageDescription' =>
      'A program running within Emacs would like to access devices on your ' \
      'local network.',

    # Core Location's macOS headers require this together with
    # NSLocationWhenInUseUsageDescription when requesting "always" access on
    # macOS 10.15 and later, despite current plist documentation presenting it
    # as an iOS key.
    'NSLocationAlwaysAndWhenInUseUsageDescription' =>
      'A program running within Emacs would like to access your location ' \
      'even when Emacs is not active.',

    # General macOS location access, required since macOS 10.14.
    'NSLocationUsageDescription' =>
      'A program running within Emacs would like to access your location.',

    # Core Location's macOS headers require this for when-in-use authorization
    # on macOS 10.15 and later, in addition to the general macOS location key.
    'NSLocationWhenInUseUsageDescription' =>
      'A program running within Emacs would like to access your location ' \
      'while Emacs is in use.',

    # Microphone access, required since macOS 10.14.
    'NSMicrophoneUsageDescription' =>
      'A program running within Emacs would like to use your microphone.',

    # Motion data access, required since macOS 10.15.
    'NSMotionUsageDescription' =>
      'A program running within Emacs would like to access motion data.',

    # Protected network-volume access, available since macOS 10.15.
    'NSNetworkVolumesUsageDescription' =>
      'A program running within Emacs would like to access files on network ' \
      'volumes.',

    # Add-only Photo Library access is available on macOS 11 and later. Apple
    # omits macOS from this plist key's current metadata, but the macOS Photos
    # framework exposes the matching PHAccessLevelAddOnly authorization mode.
    'NSPhotoLibraryAddUsageDescription' =>
      'A program running within Emacs would like to add items to your photo ' \
      'library.',

    # Read/write Photo Library access, required since macOS 10.14.
    'NSPhotoLibraryUsageDescription' =>
      'A program running within Emacs would like to access your photo library.',

    # Full reminders access on macOS 14 and later.
    'NSRemindersFullAccessUsageDescription' =>
      'A program running within Emacs would like to read and modify your ' \
      'reminders.',

    # Legacy reminders access for macOS 10.14 through macOS 13. Keep this with
    # NSRemindersFullAccessUsageDescription for older supported systems.
    'NSRemindersUsageDescription' =>
      'A program running within Emacs would like to access your reminders.',

    # Protected removable-volume access, available since macOS 10.15.
    'NSRemovableVolumesUsageDescription' =>
      'A program running within Emacs would like to access files on ' \
      'removable volumes.',

    # Speech recognition access, required since macOS 10.15.
    'NSSpeechRecognitionUsageDescription' =>
      'A program running within Emacs would like to use speech recognition.',

    # APIs that manipulate system configuration, required since macOS 10.14.
    'NSSystemAdministrationUsageDescription' =>
      'A program running within Emacs would like to modify system ' \
      'configuration.'
  }.freeze

  attr_reader :app, :editor

  def initialize(app, editor:)
    @app = app
    @editor = editor
  end

  def embed
    unless File.file?(info_plist)
      raise Error, "Info.plist not found: #{info_plist}"
    end

    DESCRIPTIONS.each do |key, description|
      editor.replace_string(info_plist, key, description)
    end
  end

  def info_plist
    @info_plist ||= File.join(app, 'Contents', 'Info.plist')
  end

  # Adapts the macOS plutil command to the small interface used here, while
  # allowing the build script to retain its command logging and error handling.
  class PlutilEditor
    def initialize(&runner)
      raise ArgumentError, 'command runner is required' unless runner

      @runner = runner
    end

    def replace_string(filename, key, value)
      @runner.call('plutil', '-replace', key, '-string', value, filename)
    end
  end
end
